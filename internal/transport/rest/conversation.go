package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"Hackathon/internal/app/grpclient"
	"Hackathon/internal/config"
	"Hackathon/internal/service"
	"Hackathon/internal/transport/dto"
	"Hackathon/internal/transport/response"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

const (
	conversationFileKey = "file"
	namePathValueKey    = "name"
	idPathValueKey      = "id"
)

type ConversationHandler struct {
	validate    *validator.Validate
	minioClient *minio.Client
	minioCfg    *config.MinioConfig
	filesCh     chan<- grpclient.FileRequest
	service     service.ConversationService
}

func NewConversationHandler(
	validate *validator.Validate,
	minioClient *minio.Client,
	minioCfg *config.MinioConfig,
	filesCh chan<- grpclient.FileRequest,
	service service.ConversationService,
) *ConversationHandler {
	return &ConversationHandler{
		validate:    validate,
		minioClient: minioClient,
		minioCfg:    minioCfg,
		filesCh:     filesCh,
		service:     service,
	}
}

// LoadConversationText docs
//
//	@Summary		Загрузка текста разговора
//	@Tags			conversation
//	@Description	Принимает текст разговора
//	@ID				load-conversation-text
//
// @Param input body dto.ConversationRequest true "текст сообщения"
//
//	@Produce		json
//	@Success		200		{object}	response.Body
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/text [post]
func (h *ConversationHandler) LoadConversationText(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	var request dto.ConversationRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(request)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	log.Printf("Got conversation text: %s", request.Text)
	response.OKMessage(w, request.Text)
}

// LoadConversationFile docs
//
//	@Summary		Загрузка аудиофайла разговора
//	@Tags			conversation
//	@Description	Принимает аудиофайл разговора
//	@ID				load-conversation-file
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200		{object}	response.Body
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/file [post]
func (h *ConversationHandler) LoadConversationFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	file, header, err := parseMultipartForm(r, conversationFileKey)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	contentType := header.Header.Get(contentTypeKey)
	extension, err := getFileExtension(contentType)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	bucketName := h.minioCfg.ConversationBucket
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	// Upload the test file with PutObject
	info, err := h.minioClient.PutObject(
		context.Background(),
		bucketName,
		header.Filename,
		bytes.NewReader(fileBytes),
		int64(len(fileBytes)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	log.Printf("Successfully uploaded %s of size %d\n", header.Filename, info.Size)

	msg := fmt.Sprintf("file extension: %s", extension)
	response.OKMessage(w, msg)
}

// GetConversationFile docs
//
//	@Summary		Получение аудиофайла разговора по его названию
//	@Tags			conversation
//	@Description	Получение аудиофайла разговора по его названию. Возвращает файл
//	@ID				get-conversation-file
//	@Produce		multipart/form-data
//	@Success		200		{file}		formData
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/file/{name} [get]
func (h *ConversationHandler) GetConversationFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	name := r.PathValue(namePathValueKey)
	if name == "" {
		response.BadRequest(w, "path value 'name' is empty")
		return
	}

	object, err := h.minioClient.GetObject(
		context.Background(),
		h.minioCfg.ConversationBucket,
		name,
		minio.GetObjectOptions{},
	)
	defer object.Close()
	log.Printf("%+v\n", object)

	if err != nil {
		log.Printf("Error while getting object: %v\n", err)

		response.InternalServerError(w)
		return
	}

	http.ServeContent(w, r, "file.mpeg", time.Now(), object)
}

// SendFileToAI docs
//
//	@Summary		Отправка аудиофайла нейронке
//	@Tags			conversation
//	@Description	Отправка аудиофайла нейронке
//	@ID				send-file-ai
//	@Produce		json
//	@Success		200		{object}	response.Body
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/file/send_ai/{name} [post]
func (h *ConversationHandler) SendFileToAI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	name := r.PathValue(namePathValueKey)
	if name == "" {
		response.BadRequest(w, "path value 'name' is empty")
		return
	}

	object, err := h.minioClient.GetObject(
		context.Background(),
		h.minioCfg.ConversationBucket,
		name,
		minio.GetObjectOptions{},
	)
	defer object.Close()
	log.Printf("%+v\n", object)

	if err != nil {
		log.Printf("Error while getting object: %v\n", err)
		response.InternalServerError(w)
		return
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(object)
	if err != nil {
		log.Printf("Error while read bytes from object: %v\n", err)
		response.InternalServerError(w)
		return
	}

	h.filesCh <- grpclient.FileRequest{
		UUID:      uuid.New().String(),
		FileName:  name,
		FileBytes: buf.Bytes(),
	}

	response.OKMessage(w, "file sent to AI successfully")
}

// GetRecords docs
//
//	@Summary		Получение записей
//	@Tags			conversation
//	@Description	Отправка записей
//	@ID				send-file-ai
//	@Produce		json
//	@Success		200		{object}	[]domain.Record
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/records [get]
func (h *ConversationHandler) GetRecords(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.GetRecords()
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	respBytes, err := json.Marshal(records)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, respBytes)
}

// GetRecord docs
//
//	@Summary		Получение записи
//	@Tags			conversation
//	@Description	Получение записи
//	@ID				send-file-ai
//	@Produce		json
//	@Success		200		{object}	domain.Record
//	@Failure		400	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Failure		default	{object}	response.Body
//	@Router			/conversation/records/{id} [get]
func (h *ConversationHandler) GetRecord(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue(idPathValueKey)
	if rawID == "" {
		response.BadRequest(w, "path value 'id' is empty")
		return
	}
	id, err := strconv.Atoi(rawID)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	record, err := h.service.GetRecord(id)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	respBytes, err := json.Marshal(record)
	if err != nil {
		log.Println(err)
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, respBytes)
}
