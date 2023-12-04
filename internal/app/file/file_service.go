package file

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"go-file-server/internal/entity"
	"io"
	"mime/multipart"
	"os"
	"path"
)

var dirPath = "./files"

type FileService struct {
	db *sql.DB
}

func NewFileService(db *sql.DB) *FileService {
	return &FileService{db: db}
}

func (s *FileService) GetAllFiles(ctx *gin.Context) ([]entity.FileInfoDto, error) {
	var files []entity.FileInfoDto

	var query = `SELECT id, file_name FROM "public"."File";`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var fileData entity.FileInfoDto
		err := rows.Scan(&fileData.ID, &fileData.Name)
		if err != nil {
			return nil, err
		}

		files = append(files, fileData)
	}

	return files, nil
}

func (s *FileService) UploadFile(header *multipart.FileHeader) error {
	file, err := header.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	dir, err := os.Getwd()
	var dirpath = path.Join(dir, "files", header.Filename)

	out, err := os.Create(dirpath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return err
	}

	var query = `INSERT INTO "public"."File"(file_name, file_path) VALUES ($1, $2)`
	_, err = s.db.Exec(query, header.Filename, dirPath+"/"+header.Filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileService) ReturnFile(id int) (string, string, error) {
	var fileName string
	var filePath string

	var query = `SELECT file_name, file_path FROM "public"."File" WHERE id = $1`
	row := s.db.QueryRow(query, id)
	err := row.Scan(&fileName, &filePath)
	if err != nil {
		return "", "", err
	}

	return fileName, filePath, nil
}
