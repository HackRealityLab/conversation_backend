package converter

import (
	"Hackathon/internal/domain"
	"Hackathon/internal/transport/dto"
)

func RecordDomainToDto(d domain.Record) dto.Record {
	return dto.Record{
		ID:          d.ID,
		Text:        d.Text,
		AudioName:   d.AudioName,
		CreatedAt:   d.CreatedAt,
		GoodPercent: d.GoodPercent,
		BadPercent:  d.BadPercent,
	}
}

func RecordSliceDomainToDto(d []domain.Record) []dto.Record {
	res := make([]dto.Record, len(d))

	for _, r := range d {
		res = append(res, RecordDomainToDto(r))
	}

	return res
}
