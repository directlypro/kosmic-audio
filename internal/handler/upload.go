package handler

import (
	"encoding/json"
	"kosmic-audio/internal/worker"
	"log/slog"
	"mime/multipart"
	"net/http"
)

type uploadHandler struct {
	pool    *worker.Pool
	storage *storage.MinioClient
}

func (h *uploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("HTTP Upload Handler")

	r.Body = http.MaxBytesReader(w, r.Body, 200<<20)

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // we can create a custom response
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			slog.Error("Error closing file", "err", err)
		}
	}(file)

	jobID := uuid.NewString()
	rawKey := "raw/" + jobID + "/" + header.Filename

	// Store raw file to MinIO before queuing
	if err := h.storage.Upload(r.Context(), rawKey, file, header.Size); err != nil {
		http.Error(w, "storage error", http.StatusInternalServerError)
		return
	}

	job := worker.Job{
		ID:           jobID,
		OriginalKey:  rawKey,
		OriginalName: header.Filename,
	}

	if err := h.pool.Submit(job); err != nil {
		http.Error(w, "storage service busy", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Content-Disposition", "attachment; filename=\""+header.Filename+"\"")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"job_id": jobID})
}
