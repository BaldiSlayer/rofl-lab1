package vdatabase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BaldiSlayer/rofl-lab1/internal/app/backend/models"
	"github.com/philippgille/chromem-go"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
)

const (
	// embeddingModel модель для эмбеддингов
	embeddingModel = "nomic-embed-text"
	// similarCount количество похожих объектов, которые мы хотим получить
	similarCount = 1

	dirPath = "./db"
)

// CMem - реализация векторной базы данных с помощью chromem-go
type CMem struct {
	collection *chromem.Collection
}

func NewCMem(baseURL string, contextFile string) (*CMem, error) {
	var cm CMem

	ctx := context.Background()

	slog.Info("Setting up chromem-go...")

	slog.Info("Sleeep...")

	time.Sleep(20 * time.Second)

	// удаляю папку ./db, чтобы каждый раз строить эмдеддинги, мне лень уже думать, как это сделать лучше
	if _, err := os.Stat(dirPath); err == nil {
		if err := os.RemoveAll(dirPath); err != nil {
			return nil, fmt.Errorf("failed to delete folder: %w", err)
		}
	}

	db, err := chromem.NewPersistentDB(dirPath, false)
	if err != nil {
		return nil, err
	}

	cm.collection, err = db.GetOrCreateCollection(
		"RoFL",
		nil,
		chromem.NewEmbeddingFuncOllama(embeddingModel, baseURL+"/api"),
	)
	if err != nil {
		return nil, err
	}

	// Add docs to the collection, if the collection was just created (and not
	// loaded from persistent storage).
	var docs []chromem.Document
	if cm.collection.Count() == 0 {
		f, err := os.Open(contextFile)
		if err != nil {
			return nil, fmt.Errorf("could not open %s: %w", contextFile, err)
		}
		defer f.Close()

		d := json.NewDecoder(f)
		for i := 1; ; i++ {
			var qaPair models.QAPair

			err := d.Decode(&qaPair)
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}

			// The embeddings model we use in this example ("nomic-embed-text")
			// fare better with a prefix to differentiate between document and query.
			content := "search_document: " + qaPair.Question

			docs = append(docs, chromem.Document{
				ID:       strconv.Itoa(i),
				Metadata: map[string]string{"answer": qaPair.Answer},
				Content:  content,
			})
		}

		slog.Info("Adding documents to chromem-go, including creating their embeddings via Ollama API...")
		err = cm.collection.AddDocuments(ctx, docs, runtime.NumCPU())
		if err != nil {
			return nil, err
		}
	} else {
		slog.Info("Not reading JSON lines because collection was loaded from persistent storage.")
	}

	slog.Info("Vector database was successfully initialized and filled")

	return &cm, nil
}

func (cm *CMem) GetSimilar(question string) ([]models.QAPair, error) {
	query := "search_query: " + question

	docRes, err := cm.collection.Query(
		context.Background(),
		query,
		similarCount,
		nil,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	contexts := make([]models.QAPair, 0, similarCount)

	for i := 0; i < similarCount; i++ {
		contexts = append(contexts, models.QAPair{
			Question: docRes[i].Content,
			Answer:   docRes[i].Metadata["answer"],
		})
	}

	return contexts, nil
}
