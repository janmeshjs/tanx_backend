package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

func CreatePriceAlert(w http.ResponseWriter, r *http.Request) {
	var alert PriceAlert
	if err := json.NewDecoder(r.Body).Decode(&alert); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	//idString := strconv.FormatUint(uint64(alert.ID), 10)

	currentPrice, err := GetCoinPrice(alert.CoinID)
	if err != nil {
		http.Error(w, "Failed to fetch current price "+err.Error(), http.StatusInternalServerError)
		return
	}

	if currentPrice <= alert.TargetPrice {
		if err := SendRabbitMQMessage(alert); err != nil {
			http.Error(w, "Failed to send RabbitMQ message", http.StatusInternalServerError)
			return
		}

		if err := SaveAlertToDB(&alert); err != nil {
			http.Error(w, "Failed to save alert to the database", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func DeletePriceAlert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	alertID := params["id"]

	var alert PriceAlert
	if err := db.Where("id = ?", alertID).First(&alert).Error; err != nil {
		http.Error(w, "Price alert not found", http.StatusNotFound)
		return
	}

	db.Delete(&alert)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Price alert with ID %s deleted", alertID)
}

func GetPriceAlerts(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	status := r.URL.Query().Get("status")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	type PriceAlertResponse struct {
		Alerts []PriceAlert `json:"alerts"`
		Status string       `json:"status"`
	}

	var alerts []PriceAlert

	query := db
	if status != "" {
		query = query.Where("status = ?", status)
	}

	totalAlertsCount := 0
	query.Find(&alerts).Count(&totalAlertsCount)

	query = query.Offset((page - 1) * 10).Limit(10).Find(&alerts)

	response := PriceAlertResponse{
		Alerts: alerts,
		Status: fmt.Sprintf("Page %d of %d", page, (totalAlertsCount+9)/10),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func SendRabbitMQMessage(alert PriceAlert) error {
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(
		"alerts",
		"fanout",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare RabbitMQ exchange: %w", err)
	}

	messageBody, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to encode message body: %w", err)
	}

	err = channel.Publish(
		"alerts", // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}

	return nil
}

func SaveAlertToDB(alert *PriceAlert) error {
	if err := db.Create(alert).Error; err != nil {
		return fmt.Errorf("failed to insert alert into database: %w", err)
	}

	return nil
}
