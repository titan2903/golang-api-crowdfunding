package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)


type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	UploadCampaignImage(input UploadCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID int) ([]Campaign, error) {
	//! melakukan switching percabangan
	if UserID != 0 {
		campaigns, err := s.repository.FindByUserID(UserID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	} else {
		campaigns, err := s.repository.FindAll()
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)

	if err != nil {
		return campaign, err
	}

	return campaign, nil
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.GoalAmount = input.GoalAmount
	campaign.Perks = input.Perks
	campaign.UserID = input.User.ID

	stringSlug := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(stringSlug)

	//! Proses pembuatan slug secara otomatis

	newCampaign, err := s.repository.SaveCampaign(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func(s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}

	if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not owner of campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updated, err := s.repository.UpdateCampaign(campaign)
	if err != nil {
		return updated, err
	}

	return updated, nil
}

func(s *service) UploadCampaignImage(input UploadCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return CampaignImage{}, err
	}

	if campaign.UserID != input.User.ID {
		return CampaignImage{}, errors.New("not owner of campaign")
	}
	
	isPrimary := 0
	if input.IsPrimary {
		isPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.UploadCampaignImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}

	return newCampaignImage, err
}