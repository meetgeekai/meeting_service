package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	gocommon "github.com/meetgeekai/go-common/utils"
	models "github.com/meetgeekai/meeting_service/internal/models"
	"github.com/meetgeekai/meeting_service/internal/utils"

	dbmodels "github.com/meetgeekai/meeting_service/internal/database/models"
)

// Keep 2 connection pools; one for slow queries and one for fast queries
type MySQLRepository struct {
	dbFast *sql.DB
	dbSlow *sql.DB
}

func NewMySQLRepository() *MySQLRepository {
	dbFast, err := sql.Open("mysql", utils.PrepareMySQLDSN())
	if err != nil {
		panic(err)
	}

	dbSlow, err := sql.Open("mysql", utils.PrepareMySQLDSN())
	if err != nil {
		panic(err)
	}

	maxConnLifetime := gocommon.GetEnv[int]("MAX_CONNECTION_LIFETIME_SECONDS")
	maxOpenConns := gocommon.GetEnv[int]("MAX_CONNECTION_POOL_SIZE")
	slowQueryPercentage := gocommon.GetEnv[int]("SLOW_QUERY_PERCENTAGE")

	maxSlowQueryConns := (maxOpenConns * slowQueryPercentage) / 100
	maxFastQueryConns := maxOpenConns - maxSlowQueryConns

	dbFast.SetConnMaxLifetime(time.Duration(maxConnLifetime) * time.Second)
	dbFast.SetMaxOpenConns(maxFastQueryConns)
	dbFast.SetMaxIdleConns(maxFastQueryConns)

	dbSlow.SetConnMaxLifetime(time.Duration(maxConnLifetime) * time.Second)
	dbSlow.SetMaxOpenConns(maxSlowQueryConns)
	dbSlow.SetMaxIdleConns(maxSlowQueryConns)

	return &MySQLRepository{
		dbFast: dbFast,
		dbSlow: dbSlow,
	}
}

func (r *MySQLRepository) GetUserForUpcomingMeetings(ctx context.Context, userID uint32) (*models.UpcomingMeetingsOwner, error) {
	queries := dbmodels.New(r.dbFast)

	row, err := queries.GetUserForUpcomingMeetings(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &models.UpcomingMeetingsOwner{
		UUID:  row.Uuid,
		Name:  row.Name,
		Email: row.Email,
	}, nil
}

func (r *MySQLRepository) GetConversationTemplateNames(ctx context.Context, ids []int64, userID uint32) (map[int64]string, error) {
	queries := dbmodels.New(r.dbFast)

	uint32IDs := make([]uint32, len(ids))
	for i, id := range ids {
		uint32IDs[i] = uint32(id)
	}

	rows, err := queries.GetConversationTemplateNames(ctx, dbmodels.GetConversationTemplateNamesParams{
		Ids:    uint32IDs,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	result := make(map[int64]string, len(rows))
	for _, row := range rows {
		result[int64(row.ID)] = row.Name
	}
	return result, nil
}

func (r *MySQLRepository) GetConnectedCalendarVendors(ctx context.Context, userUUID string) (models.ConnectedCalendars, error) {
	queries := dbmodels.New(r.dbFast)

	vendors, err := queries.GetConnectedCalendarVendors(ctx, userUUID)
	if err != nil {
		return models.ConnectedCalendars{}, err
	}

	var result models.ConnectedCalendars
	for _, vendor := range vendors {
		provider, ok := models.ParseCalendarProvider(vendor)
		if !ok {
			continue
		}
		switch provider {
		case models.CalendarProviderGoogle:
			result.Google = true
		case models.CalendarProviderMicrosoft:
			result.Microsoft = true
		}
	}
	return result, nil
}

func (r *MySQLRepository) GetAvailableTranscriptionLanguages(ctx context.Context) ([]models.TranscriptionLanguage, error) {
	queries := dbmodels.New(r.dbFast)

	rows, err := queries.GetAvailableTranscriptionLanguages(ctx)
	if err != nil {
		return nil, err
	}

	languages := make([]models.TranscriptionLanguage, 0, len(rows))
	for _, row := range rows {
		languages = append(languages, models.TranscriptionLanguage{
			ID:       int64(row.ID),
			Code:     string(row.Code),
			Value:    row.Value,
			Country:  row.Country,
			Language: row.Language,
			CustomDictionary: row.BoostMeetingAssistant.Valid ||
				row.BoostMeetingTitle.Valid ||
				row.BoostMeetingParticipants.Valid ||
				row.BoostMeetingParticipantsActual.Valid ||
				row.BoostTemplateCustomWords.Valid ||
				row.BoostTemplateKeywordBasedHighlights.Valid,
		})
	}

	return languages, nil
}
