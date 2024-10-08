// This file contains the SceneManager implementation, which is responsible for interacting with the MongoDB scene collection.
// The SceneManager struct contains a pointer to the nerfdb.scenes MongoDB collection and a logger. It provides methods to set and
// get scene data from the database. Interaction with scenes is almost always by ID, as the ID will (almost always) be unique.

package scene

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/NeRF-or-Nothing/go-web-server/internal/log"
)

// Custom errors
var (
	// ErrSceneNotFound is returned when a requested scene is not found in the database.
	ErrSceneNotFound = errors.New("scene not found")
	// ErrVideoNotFound is returned when a requested video is not found in the database.
	ErrVideoNotFound = errors.New("video not found")
	// ErrSfmNotFound is returned when a requested sfm is not found in the database.
	ErrSfmNotFound = errors.New("sfm not found")
	// ErrNerfNotFound is returned when a requested nerf is not found in the database.
	ErrNerfNotFound = errors.New("nerf not found")
	// ErrTrainingConfigNotFound is returned when a requested training config is not found in the database.
	ErrTrainingConfigNotFound = errors.New("training config not found")
)

type SceneManager struct {
	collection *mongo.Collection
	logger     *log.Logger
}

// NewSceneManager creates a new SceneManager with the given MongoDB client and logger.
func NewSceneManager(client *mongo.Client, logger *log.Logger, unittest bool) *SceneManager {
	return &SceneManager{
		collection: client.Database("nerfdb").Collection("scenes"),
		logger:     logger,
	}
}

// SetTrainingConfig sets the TrainingConfig data in the database by the scene ID.
func (sm *SceneManager) SetTrainingConfig(ctx context.Context, id primitive.ObjectID, config *TrainingConfig) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"config": config}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// SetScene sets the Scene data in the database by the scene ID.
func (sm *SceneManager) SetScene(ctx context.Context, id primitive.ObjectID, scene *Scene) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": scene},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// SetVideo sets the Video data in the database by the scene ID.
func (sm *SceneManager) SetVideo(ctx context.Context, id primitive.ObjectID, vid *Video) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"video": vid}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// SetSfm sets the Sfm data in the database by the scene ID.
func (sm *SceneManager) SetSfm(ctx context.Context, id primitive.ObjectID, sfm *Sfm) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"sfm": sfm}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// SetNerf sets the Nerf data in the database by the scene ID.
func (sm *SceneManager) SetNerf(ctx context.Context, id primitive.ObjectID, nerf *Nerf) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"nerf": nerf}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// SetSceneName sets the name of the scene in the database by its ID.
func (sm *SceneManager) SetSceneName(ctx context.Context, id primitive.ObjectID, name string) error {
	result, err := sm.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"name": name}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 && result.UpsertedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}

// GetSceneName retrieves the name of the scene from the database by its ID.
func (sm *SceneManager) GetSceneName(ctx context.Context, id primitive.ObjectID) (string, error) {
	var result struct {
		Name string `bson:"name"`
	}
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", ErrSceneNotFound
		}
		return "", err
	}
	return result.Name, nil
}

// GetTrainingConfig retrieves the TrainingConfig data from the database by its ID.
func (sm *SceneManager) GetTrainingConfig(ctx context.Context, id primitive.ObjectID) (*TrainingConfig, error) {
	var result struct {
		Config *TrainingConfig `bson:"config"`
	}
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	if result.Config == nil {
		return nil, ErrTrainingConfigNotFound
	}
	return result.Config, nil
}

// GetScene retrieves the Scene data from the database by its ID.
func (sm *SceneManager) GetScene(ctx context.Context, id primitive.ObjectID) (*Scene, error) {
	var scene Scene
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&scene)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	return &scene, nil
}

// GetVideo retrieves the Video data from the database by its ID.
func (sm *SceneManager) GetVideo(ctx context.Context, id primitive.ObjectID) (*Video, error) {
	var result struct {
		Video *Video `bson:"video"`
	}
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	if result.Video == nil {
		return nil, ErrVideoNotFound
	}
	return result.Video, nil
}

// GetSfm retrieves the Sfm data from the database by its ID.
func (sm *SceneManager) GetSfm(ctx context.Context, id primitive.ObjectID) (*Sfm, error) {
	var result struct {
		Sfm *Sfm `bson:"sfm"`
	}
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	if result.Sfm == nil {
		return nil, ErrSfmNotFound
	}
	return result.Sfm, nil
}

// GetNerf retrieves the Nerf data from the database by its ID.
func (sm *SceneManager) GetNerf(ctx context.Context, id primitive.ObjectID) (*Nerf, error) {
	var result struct {
		Nerf *Nerf `bson:"nerf"`
	}
	err := sm.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrSceneNotFound
		}
		return nil, err
	}
	if result.Nerf == nil {
		return nil, ErrNerfNotFound
	}
	return result.Nerf, nil
}

// DeleteScene deletes a scene from the database by its ID.
func (sm *SceneManager) DeleteScene(ctx context.Context, id primitive.ObjectID) error {
	result, err := sm.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrSceneNotFound
	}
	return nil
}
