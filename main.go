package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type SimrsRequest struct {
	NoPendaftaran string `json:"no_pendaftaran"`
	NoRM          string `json:"no_rm"`
	NoOrder       string `json:"no_order"`
	NamaPasien    string `json:"nama_pasien"`
	TempatLahir   string `json:"tempat_lahir"`
	TglLahir      string `json:"tgl_lahir"`
	JK            string `json:"jk"`
	Alamat        string `json:"alamat"`
	Phone         string `json:"phone,omitempty"`
	NIK           string `json:"nik"`
	IdJenisPasien string `json:"id_jenis_pasien"`
	JenisPasien   string `json:"jenis_pasien,omitempty"`
	IdPenjamin    string `json:"id_penjamin"`
	Penjamin      string `json:"penjamin,omitempty"`
	RujukanAsal   string `json:"rujukan_asal"`
	DetailRujukan []struct {
		IdDokter    string `json:"id_dokter"`
		NamaDokter  string `json:"nama_dokter"`
		IdWard      string `json:"id_ward"`
		Ward        string `json:"ward"`
		IdFasilitas string `json:"id_fasilitas"`
		Fasilitas   string `json:"fasilitas"`
	} `json:"detail_rujukan"`
	Cito              string   `json:"cito"`
	Diagnose          string   `json:"diagnose"`
	RegistrationDate  string   `json:"registration_date"`
	Email             string   `json:"email"`
	Region            string   `json:"region"`
	CodeClinic        string   `json:"code_clinic"`
	ClinicName        string   `json:"clinic_name"`
	OrderNote         string   `json:"order_note"`
	KategoriKunjungan string   `json:"kategori_kunjungan"`
	ICD10             []string `json:"icd10"`
	Order             []struct {
		IdTest   string `json:"id_test"`
		NamaTest string `json:"nama_test"`
	} `json:"order"`
}

type LISRequest struct {
	NoPendaftaran string `json:"no_pendaftaran"`
	NoRM          string `json:"no_rm"`
	NoOrder       string `json:"no_order"`
	NamaPasien    string `json:"nama_pasien"`
	TempatLahir   string `json:"tempat_lahir"`
	TglLahir      string `json:"tgl_lahir"`
	JK            string `json:"jk"`
	Alamat        string `json:"alamat"`
	Phone         string `json:"phone,omitempty"`
	NIK           string `json:"nik"`
	IdJenisPasien string `json:"id_jenis_pasien"`
	JenisPasien   string `json:"jenis_pasien,omitempty"`
	IdPenjamin    string `json:"id_penjamin"`
	Penjamin      string `json:"penjamin,omitempty"`
	RujukanAsal   string `json:"rujukan_asal"`
	DetailRujukan []struct {
		IdDokter    string `json:"id_dokter"`
		NamaDokter  string `json:"nama_dokter"`
		IdWard      string `json:"id_ward"`
		Ward        string `json:"ward"`
		IdFasilitas string `json:"id_fasilitas"`
		Fasilitas   string `json:"fasilitas"`
	} `json:"detail_rujukan"`
	Cito              bool     `json:"cito"`
	Diagnose          string   `json:"diagnose"`
	Email             string   `json:"email"`
	Region            string   `json:"region"`
	CodeClinic        string   `json:"code_clinic"`
	ClinicName        string   `json:"clinic_name"`
	OrderNote         string   `json:"order_note"`
	KategoriKunjungan string   `json:"kategori_kunjungan"`
	ICD10             []string `json:"icd10"`
	Order             []struct {
		IdTest   string `json:"id_test"`
		NamaTest string `json:"nama_test"`
	} `json:"order"`
}

func transformToLIS(simrs *SimrsRequest) *LISRequest {
	var idJenisPasien, jenisPasien, idPenjamin, penjamin string

	if simrs.IdJenisPasien == "UMUM" {
		idJenisPasien = "1"
		jenisPasien = "UMUM"
		idPenjamin = "UMUM"
		penjamin = " "
	} else {
		idJenisPasien = "2"
		jenisPasien = "ASURANSI"
		idPenjamin = simrs.IdJenisPasien
		penjamin = simrs.JenisPasien
	}

	// Transform Cito to bool
	var cito bool
	if simrs.Cito == "1" || simrs.Cito == "true" {
		cito = true
	} else {
		cito = false
	}

	// Map SimrsRequest to LISRequest
	return &LISRequest{
		NoPendaftaran:     simrs.NoPendaftaran,
		NoRM:              simrs.NoRM,
		NoOrder:           simrs.NoOrder,
		NamaPasien:        simrs.NamaPasien,
		TempatLahir:       simrs.TempatLahir,
		TglLahir:          simrs.TglLahir,
		JK:                simrs.JK,
		Alamat:            simrs.Alamat,
		Phone:             simrs.Phone,
		NIK:               simrs.NIK,
		IdJenisPasien:     idJenisPasien,
		JenisPasien:       jenisPasien,
		IdPenjamin:        idPenjamin,
		Penjamin:          penjamin,
		RujukanAsal:       simrs.RujukanAsal,
		DetailRujukan:     simrs.DetailRujukan,
		Cito:              cito,
		Diagnose:          simrs.Diagnose,
		Email:             simrs.Email,
		Region:            simrs.Region,
		CodeClinic:        simrs.CodeClinic,
		ClinicName:        simrs.ClinicName,
		OrderNote:         simrs.OrderNote,
		KategoriKunjungan: simrs.KategoriKunjungan,
		ICD10:             simrs.ICD10,
		Order:             simrs.Order,
	}
}

type SendToLisBridgingIn struct {
	request *LISRequest
	xSign   string
	xCons   string
}

func sendToLISBridging(payload SendToLisBridgingIn) (map[string]interface{}, error) {
    url := os.Getenv("LIS_BRIDGING")

    jsonData, err := json.Marshal(payload.request)
    var respBody map[string]interface{}

    if err != nil {
        return respBody, err
    }

    log.Printf("Sending data to LIS bridging: %s", jsonData)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return respBody, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-sign", payload.xSign)
    req.Header.Set("x-cons", payload.xCons)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return respBody, err
    }
    defer resp.Body.Close()

    // Log the response from the API
    err = json.NewDecoder(resp.Body).Decode(&respBody)
    if err != nil {
        return respBody, err
    }
    log.Printf("Response from LIS bridging: %s", resp.Status)
    log.Printf("Response body: %v", respBody)

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return respBody, fmt.Errorf("failed to send data to LIS bridging, status code: %d", resp.StatusCode)
    }

    return respBody, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var simrsReq SimrsRequest

	err := json.NewDecoder(r.Body).Decode(&simrsReq)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Transform SIMRS DTO to LIS DTO
	lisReq := transformToLIS(&simrsReq)

	// Send the transformed data to LIS bridging
	payload := SendToLisBridgingIn{
		request: lisReq,
		xSign:   r.Header.Get("x-sign"),
		xCons:   r.Header.Get("x-cons"),
	}

	respBody, err := sendToLISBridging(payload)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending data to LIS bridging: %v", err), http.StatusInternalServerError)
		log.Printf("Error sending data to LIS bridging: %v", err)
		return
	}

	// Respond with the transformed data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
	log.Printf("\n------------------------------------------------------------------------------------------------------------------\n")

	// Log the response status
	log.Printf("Response status:\n %d", http.StatusOK)
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Post("/api/v1/saveOrder", handleRequest)

	httpServer := &http.Server{
		Addr:    ":8111",
		Handler: r,
	}

	log.Printf("Starting server on port %s\n", "8111")
	log.Fatal(httpServer.ListenAndServe())
}
