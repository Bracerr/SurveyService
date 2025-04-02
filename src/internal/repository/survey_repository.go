package repository

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"survey-project/src/internal/apperrors"
	"survey-project/src/internal/domain"
)

type SurveyRepository struct {
	collection *mongo.Collection
}

func NewSurveyRepository(collection *mongo.Collection) *SurveyRepository {
	return &SurveyRepository{collection: collection}
}

func (r *SurveyRepository) Create(survey *domain.Survey) error {
	ctx := context.Background()

	id := primitive.NewObjectID().Hex()
	survey.ID = id

	for i := range survey.Questions {
		survey.Questions[i].ID = primitive.NewObjectID().Hex()
	}

	survey.CreatedAt = time.Now()
	survey.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, bson.M{
		"_id":          id,
		"id":           id, 
		"title":        survey.Title,
		"description":  survey.Description,
		"created_by":   survey.CreatedBy,
		"questions":    survey.Questions,
		"is_active":    survey.IsActive,
		"is_anonymous": survey.IsAnonymous,
		"require_info": survey.RequireInfo,
		"start_date":   survey.StartDate,
		"end_date":     survey.EndDate,
		"created_at":   survey.CreatedAt,
		"updated_at":   survey.UpdatedAt,
	})

	return err
}

func (r *SurveyRepository) GetByID(id string) (*domain.Survey, error) {
	ctx := context.Background()

	var survey domain.Survey
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&survey)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.ErrSurveyNotFound
		}
		return nil, fmt.Errorf("failed to get survey: %w", err)
	}

	survey.ID = id

	return &survey, nil
}

func (r *SurveyRepository) GetAll() ([]*domain.Survey, error) {
	ctx := context.Background()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get surveys: %w", err)
	}
	defer cursor.Close(ctx)

	var rawDocs []bson.M
	if err := cursor.All(ctx, &rawDocs); err != nil {
		return nil, fmt.Errorf("failed to decode surveys: %w", err)
	}

	surveys := make([]*domain.Survey, len(rawDocs))
	for i, doc := range rawDocs {
		var survey domain.Survey
		bytes, err := bson.Marshal(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal document: %w", err)
		}
		if err := bson.Unmarshal(bytes, &survey); err != nil {
			return nil, fmt.Errorf("failed to unmarshal survey: %w", err)
		}

		if oid, ok := doc["_id"].(string); ok {
			survey.ID = oid
		}
		surveys[i] = &survey
	}

	return surveys, nil
}

func (r *SurveyRepository) Update(survey *domain.Survey) error {
	ctx := context.Background()

	surveyID, ok := survey.ID.(string)
	if !ok {
		return fmt.Errorf("invalid survey ID type")
	}

	updateFields := bson.M{"updated_at": time.Now()}

	val := reflect.ValueOf(survey).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Name == "ID" || fieldType.Name == "CreatedAt" || fieldType.Name == "UpdatedAt" {
			continue
		}

		if !isZeroValue(field) {
			bsonTag := fieldType.Tag.Get("bson")
			if bsonTag == "" || bsonTag == "-" {
				continue
			}
			bsonName := strings.Split(bsonTag, ",")[0]

			updateFields[bsonName] = field.Interface()
		}
	}

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": surveyID},
		bson.M{"$set": updateFields},
	)
	return err
}

func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return isZeroValue(v.Elem())
	case reflect.Bool:
		return !v.Bool()
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Struct:
		if t, ok := v.Interface().(time.Time); ok {
			return t.IsZero()
		}
		for i := 0; i < v.NumField(); i++ {
			if !isZeroValue(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (r *SurveyRepository) Delete(id string) error {
	ctx := context.Background()

	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *SurveyRepository) GetByUserID(userID string) ([]*domain.Survey, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{"created_by": userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user surveys: %w", err)
	}
	defer cursor.Close(context.Background())

	var surveys []*domain.Survey
	if err := cursor.All(context.Background(), &surveys); err != nil {
		return nil, fmt.Errorf("failed to decode user surveys: %w", err)
	}

	for _, survey := range surveys {
		if id, ok := survey.ID.(string); ok {
			survey.ID = id
		}
	}

	return surveys, nil
}
