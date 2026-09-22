package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Member struct {
	ID                    int    `json:"id"`
	Name                  string `json:"name"`
	PersonalEmail         string `json:"personalEmail"`
	Branch                string `json:"branch"`
	ExpectedYearOfPassing int    `json:"expectedYearOfPassing"`
	Phone                 string `json:"phone"`
	Domain                string `json:"domain"`
	TempRtfID             string `json:"tempRtfId"`
	Status                string `json:"status"`
}

var (
	members = make(map[int]Member)
	nextID  = 1
	mu      sync.Mutex
)

func main() {

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/members", membersHandler)
	http.HandleFunc("/members/", memberHandler)

	fmt.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// GET /health
func healthHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "success",
		"message": "Member API is running",
	}

	json.NewEncoder(w).Encode(response)
}

// POST /members
// GET /members
func membersHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		createMember(w, r)

	case http.MethodGet:
		getMembers(w)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /members/{id}
// PUT /members/{id}
// DELETE /members/{id}
func memberHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	idString := strings.TrimPrefix(r.URL.Path, "/members/")

	id, err := strconv.Atoi(idString)

	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		getMember(w, id)

	case http.MethodPut:
		updateMember(w, r, id)

	case http.MethodDelete:
		deleteMember(w, id)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CREATE MEMBER
func getDomainCode(domain string) string {
	switch domain {
	case "software":
		return "SD"
	case "electrical":
		return "EL"
	case "aeromech":
		return "AM"
	default:
		return ""
	}
}
func createMember(w http.ResponseWriter, r *http.Request) {

	var member Member

	err := json.NewDecoder(r.Body).Decode(&member)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate domain
	if member.Domain != "software" &&
		member.Domain != "electrical" &&
		member.Domain != "aeromech" {

		http.Error(w, "Invalid domain", http.StatusBadRequest)
		return
	}

	// Validate status
	if member.Status != "applied" &&
		member.Status != "converted" {

		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	mu.Lock()

	member.ID = nextID
	nextID++

	domainCode := getDomainCode(member.Domain)
	year := member.ExpectedYearOfPassing % 100
	serial := member.ID

	member.TempRtfID = fmt.Sprintf("%s%02d%02d@RTF", domainCode, year, serial)

	members[member.ID] = member

	mu.Unlock()

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(member)
}

// GET ALL MEMBERS
func getMembers(w http.ResponseWriter) {

	mu.Lock()
	defer mu.Unlock()

	memberList := make([]Member, 0, len(members))

	for _, member := range members {
		memberList = append(memberList, member)
	}

	json.NewEncoder(w).Encode(memberList)
}

// GET ONE MEMBER
func getMember(w http.ResponseWriter, id int) {

	mu.Lock()
	defer mu.Unlock()

	member, exists := members[id]

	if !exists {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(member)
}

// UPDATE MEMBER
func updateMember(w http.ResponseWriter, r *http.Request, id int) {

	var updatedMember Member

	err := json.NewDecoder(r.Body).Decode(&updatedMember)

	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	_, exists := members[id]

	if !exists {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	updatedMember.ID = id

	members[id] = updatedMember

	json.NewEncoder(w).Encode(updatedMember)
}

// DELETE MEMBER
func deleteMember(w http.ResponseWriter, id int) {

	mu.Lock()
	defer mu.Unlock()

	_, exists := members[id]

	if !exists {
		http.Error(w, "Member not found", http.StatusNotFound)
		return
	}

	delete(members, id)

	response := map[string]string{
		"message": "Member deleted successfully",
	}

	json.NewEncoder(w).Encode(response)
}
