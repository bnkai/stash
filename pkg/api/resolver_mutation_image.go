package api

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) ImageUpdate(ctx context.Context, input models.ImageUpdateInput) (ret *models.Image, err error) {
	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	// Start the transaction and save the image
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		ret, err = r.imageUpdate(input, translator, repo)
		return err
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) ImagesUpdate(ctx context.Context, input []*models.ImageUpdateInput) (ret []*models.Image, err error) {
	inputMaps := getUpdateInputMaps(ctx)

	// Start the transaction and save the image
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		for i, image := range input {
			translator := changesetTranslator{
				inputMap: inputMaps[i],
			}

			thisImage, err := r.imageUpdate(*image, translator, repo)
			if err != nil {
				return err
			}

			ret = append(ret, thisImage)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *mutationResolver) imageUpdate(input models.ImageUpdateInput, translator changesetTranslator, repo models.Repository) (*models.Image, error) {
	// Populate image from the input
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return nil, err
	}

	updatedTime := time.Now()
	updatedImage := models.ImagePartial{
		ID:        imageID,
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	updatedImage.Title = translator.nullString(input.Title, "title")
	updatedImage.Rating = translator.nullInt64(input.Rating, "rating")
	updatedImage.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedImage.Organized = input.Organized

	qb := repo.Image()
	image, err := qb.Update(updatedImage)
	if err != nil {
		return nil, err
	}

	if translator.hasField("gallery_ids") {
		if err := r.updateImageGalleries(qb, imageID, input.GalleryIds); err != nil {
			return nil, err
		}
	}

	// Save the performers
	if translator.hasField("performer_ids") {
		if err := r.updateImagePerformers(qb, imageID, input.PerformerIds); err != nil {
			return nil, err
		}
	}

	// Save the tags
	if translator.hasField("tag_ids") {
		if err := r.updateImageTags(qb, imageID, input.TagIds); err != nil {
			return nil, err
		}
	}

	return image, nil
}

func (r *mutationResolver) updateImageGalleries(qb models.ImageReaderWriter, imageID int, galleryIDs []string) error {
	ids, err := utils.StringSliceToIntSlice(galleryIDs)
	if err != nil {
		return err
	}
	return qb.UpdateGalleries(imageID, ids)
}

func (r *mutationResolver) updateImagePerformers(qb models.ImageReaderWriter, imageID int, performerIDs []string) error {
	ids, err := utils.StringSliceToIntSlice(performerIDs)
	if err != nil {
		return err
	}
	return qb.UpdatePerformers(imageID, ids)
}

func (r *mutationResolver) updateImageTags(qb models.ImageReaderWriter, imageID int, tagsIDs []string) error {
	ids, err := utils.StringSliceToIntSlice(tagsIDs)
	if err != nil {
		return err
	}
	return qb.UpdateTags(imageID, ids)
}

func (r *mutationResolver) BulkImageUpdate(ctx context.Context, input models.BulkImageUpdateInput) (ret []*models.Image, err error) {
	imageIDs, err := utils.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return nil, err
	}

	// Populate image from the input
	updatedTime := time.Now()

	updatedImage := models.ImagePartial{
		UpdatedAt: &models.SQLiteTimestamp{Timestamp: updatedTime},
	}

	translator := changesetTranslator{
		inputMap: getUpdateInputMap(ctx),
	}

	updatedImage.Title = translator.nullString(input.Title, "title")
	updatedImage.Rating = translator.nullInt64(input.Rating, "rating")
	updatedImage.StudioID = translator.nullInt64FromString(input.StudioID, "studio_id")
	updatedImage.Organized = input.Organized

	// Start the transaction and save the image marker
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		for _, imageID := range imageIDs {
			updatedImage.ID = imageID

			image, err := qb.Update(updatedImage)
			if err != nil {
				return err
			}

			ret = append(ret, image)

			// Save the galleries
			if translator.hasField("gallery_ids") {
				galleryIDs, err := adjustImageGalleryIDs(qb, imageID, *input.GalleryIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateGalleries(imageID, galleryIDs); err != nil {
					return err
				}
			}

			// Save the performers
			if translator.hasField("performer_ids") {
				performerIDs, err := adjustImagePerformerIDs(qb, imageID, *input.PerformerIds)
				if err != nil {
					return err
				}

				if err := qb.UpdatePerformers(imageID, performerIDs); err != nil {
					return err
				}
			}

			// Save the tags
			if translator.hasField("tag_ids") {
				tagIDs, err := adjustImageTagIDs(qb, imageID, *input.TagIds)
				if err != nil {
					return err
				}

				if err := qb.UpdateTags(imageID, tagIDs); err != nil {
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return ret, nil
}

func adjustImageGalleryIDs(qb models.ImageReader, imageID int, ids models.BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetGalleryIDs(imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustImagePerformerIDs(qb models.ImageReader, imageID int, ids models.BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetPerformerIDs(imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func adjustImageTagIDs(qb models.ImageReader, imageID int, ids models.BulkUpdateIds) (ret []int, err error) {
	ret, err = qb.GetTagIDs(imageID)
	if err != nil {
		return nil, err
	}

	return adjustIDs(ret, ids), nil
}

func (r *mutationResolver) ImageDestroy(ctx context.Context, input models.ImageDestroyInput) (ret bool, err error) {
	imageID, err := strconv.Atoi(input.ID)
	if err != nil {
		return false, err
	}

	var image *models.Image
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		image, err = qb.Find(imageID)
		if err != nil {
			return err
		}

		if image == nil {
			return fmt.Errorf("image with id %d not found", imageID)
		}

		return qb.Destroy(imageID)
	}); err != nil {
		return false, err
	}

	// if delete generated is true, then delete the generated files
	// for the image
	if input.DeleteGenerated != nil && *input.DeleteGenerated {
		manager.DeleteGeneratedImageFiles(image)
	}

	// if delete file is true, then delete the file as well
	// if it fails, just log a message
	if input.DeleteFile != nil && *input.DeleteFile {
		manager.DeleteImageFile(image)
	}

	return true, nil
}

func (r *mutationResolver) ImagesDestroy(ctx context.Context, input models.ImagesDestroyInput) (ret bool, err error) {
	imageIDs, err := utils.StringSliceToIntSlice(input.Ids)
	if err != nil {
		return false, err
	}

	var images []*models.Image
	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		for _, imageID := range imageIDs {

			image, err := qb.Find(imageID)
			if err != nil {
				return err
			}

			if image == nil {
				return fmt.Errorf("image with id %d not found", imageID)
			}

			images = append(images, image)
			if err := qb.Destroy(imageID); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return false, err
	}

	for _, image := range images {
		// if delete generated is true, then delete the generated files
		// for the image
		if input.DeleteGenerated != nil && *input.DeleteGenerated {
			manager.DeleteGeneratedImageFiles(image)
		}

		// if delete file is true, then delete the file as well
		// if it fails, just log a message
		if input.DeleteFile != nil && *input.DeleteFile {
			manager.DeleteImageFile(image)
		}
	}

	return true, nil
}

func (r *mutationResolver) ImageIncrementO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.IncrementOCounter(imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageDecrementO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.DecrementOCounter(imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}

func (r *mutationResolver) ImageResetO(ctx context.Context, id string) (ret int, err error) {
	imageID, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	if err := r.withTxn(ctx, func(repo models.Repository) error {
		qb := repo.Image()

		ret, err = qb.ResetOCounter(imageID)
		return err
	}); err != nil {
		return 0, err
	}

	return ret, nil
}
