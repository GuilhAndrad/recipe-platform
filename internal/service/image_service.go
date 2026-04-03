package service

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
 
	"github.com/disintegration/imaging"
)
 
const (
	uploadDir    = "uploads"
	targetWidth  = 800
	targetHeight = 600
)
 
type ImageService struct{}
 
func NewImageService() *ImageService {
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(fmt.Sprintf("Não foi possível criar o diretório de uploads: %v", err))
	}
	return &ImageService{}
}
 
//recebe o arquivo, redimensiona para 800x600 e salva em disco.
func (s *ImageService) Process(fileHeader *multipart.FileHeader) (string, error) {
	if err := s.validateExtension(fileHeader.Filename); err != nil {
		return "", err
	}
 
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	defer src.Close()
 
	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return "", fmt.Errorf("arquivo não é uma imagem válida: %w", err)
	}
 
	// Fill garante 800x600 exatos sem distorcer — corta o excesso pelo centro
	resized := imaging.Fill(img, targetWidth, targetHeight, imaging.Center, imaging.Lanczos)
 
	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	savePath := filepath.Join(uploadDir, filename)
 
	if err := imaging.Save(resized, savePath); err != nil {
		return "", fmt.Errorf("erro ao salvar imagem: %w", err)
	}
 
	return savePath, nil
}
 
//remove a imagem do disco quando uma receita é excluída.
func (s *ImageService) Delete(imagePath string) error {
	if imagePath == "" {
		return nil
	}
	_ = os.Remove(imagePath)
	return nil
}
 
func (s *ImageService) validateExtension(filename string) error {
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true,
		".png": true, ".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowed[ext] {
		return fmt.Errorf("formato não suportado: %s. Use jpg, png ou webp", ext)
	}
	return nil
}