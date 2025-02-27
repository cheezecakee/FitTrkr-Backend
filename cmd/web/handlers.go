package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/cheezecakee/FitLogr/internal/database"
)

// Handler struct to hold dbQueries
type ApiConfig struct {
	DB        *database.Queries
	JWTSecret []byte
}

// Get Workouts from JSON
func (apiCfg *ApiConfig) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch workouts from DB
	workouts, err := apiCfg.DB.GetWorkoutsByID(ctx, userID)
	if err != nil {
		http.Error(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(workouts); err != nil {
		http.Error(w, "Failed to encode workouts", http.StatusInternalServerError)
	}
}

// Create Workout
func (apiCfg *ApiConfig) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse JSON request body
	var req Workout

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Insert workout into DB
	workout, err := apiCfg.DB.CreateWorkout(ctx, database.CreateWorkoutParams{
		UserID:      userID,
		Name:        req.Name,
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		log.Printf("Error creating workout: %v", err) // Log the actual error from the database
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	log.Printf("Workout created succesfully!") // Log the actual error from the database
	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// Create Exercise
func (apiCfg *ApiConfig) CreateWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req WorkoutExercise

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Convert workout_id and exercise_id to UUID
	workoutID, err := uuid.Parse(req.WorkoutID)
	if err != nil {
		http.Error(w, "Invalid workout ID format", http.StatusBadRequest)
		return
	}

	exerciseID, err := uuid.Parse(req.ExerciseID)
	if err != nil {
		http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	// Convert rest_min to seconds and add to rest_sec
	totalRestSeconds := (req.RestMin * 60) + req.RestSec

	// Create workout exercise in the database
	workoutExercise, err := apiCfg.DB.CreateWorkoutExercise(ctx, database.CreateWorkoutExerciseParams{
		WorkoutID:  workoutID,
		ExerciseID: exerciseID,
		Sets:       req.Sets,
		RepsMin:    req.RepsMin,
		RepsMax:    req.RepsMax,
		Weight:     req.Weight,
		Interval:   req.Interval,
		Rest:       totalRestSeconds,
	})
	if err != nil {
		log.Printf("Error creating workout exercise: %v", err) // Log the actual error from the database
		http.Error(w, "Failed to create workout exercise", http.StatusInternalServerError)
		return
	}

	log.Printf("Exercise added to workout succesfully!") // Log the actual error from the database

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workoutExercise)
}

func (apiCfg *ApiConfig) EditWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workoutIDStr := r.PathValue("id")

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		http.Error(w, "Invalid workout ID format", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update workout
	updatedWorkout, err := apiCfg.DB.EditWorkout(ctx, database.EditWorkoutParams{
		ID:          workoutID,
		UserID:      userID,
		Name:        StringToNullString(req.Name),
		Description: StringToNullString(req.Description),
	})
	if err != nil {
		log.Printf("%w", err)
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}

	log.Printf("Workout Exercise updated succesfully!") // Log the actual error from the database
	// Return updated workout
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorkout)
}

func (apiCfg *ApiConfig) EditWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exerciseIDStr := r.PathValue("id")

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	var req struct {
		Sets     int32   `json:"sets"`
		RepsMin  int32   `json:"reps_min"`
		RepsMax  int32   `json:"reps_max"`
		Weight   float64 `json:"weight"`
		Interval int32   `json:"interval"`
		Rest     int32   `json:"rest"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update workout exercise
	updatedExercise, err := apiCfg.DB.EditWorkoutExercise(ctx, database.EditWorkoutExerciseParams{
		ID:       exerciseID,
		Sets:     sql.NullInt32{Int32: req.Sets, Valid: true},
		RepsMin:  sql.NullInt32{Int32: req.RepsMin, Valid: true},
		RepsMax:  sql.NullInt32{Int32: req.RepsMax, Valid: true},
		Weight:   sql.NullFloat64{Float64: req.Weight, Valid: true},
		Interval: sql.NullInt32{Int32: req.Interval, Valid: true},
		Rest:     sql.NullInt32{Int32: req.Rest, Valid: true},
	})
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Failed to update workout exercise", http.StatusInternalServerError)
		return
	}

	log.Printf("Workout Exercise updated succesfully!") // Log the actual error from the database
	// Return updated exercise
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedExercise)
}

func (apiCfg *ApiConfig) DeleteWorkoutExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	exerciseIDStr := r.PathValue("id")

	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID format", http.StatusBadRequest)
		return
	}

	// Execute delete query
	err = apiCfg.DB.DeleteWorkoutExercise(ctx, exerciseID)
	if err != nil {
		http.Error(w, "Failed to delete workout exercise", http.StatusInternalServerError)
		return
	}

	log.Printf("Workout Exercise deleted succesfully!") // Log the actual error from the database
	w.WriteHeader(http.StatusNoContent)                 // 204 No Content response
}

func (apiCfg *ApiConfig) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract workout ID from URL
	workoutIDStr := r.PathValue("id")
	log.Printf("workoutID: %s", workoutIDStr)
	workoutID, err := uuid.Parse(workoutIDStr)
	if err != nil {
		http.Error(w, "Invalid workout ID format", http.StatusBadRequest)
		return
	}

	// Execute delete query
	err = apiCfg.DB.DeleteWorkout(ctx, workoutID)
	if err != nil {
		log.Printf("Error deleting workout: %v", err)
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}
	log.Printf("Workout deleted succesfully!") // Log the actual error from the database

	w.WriteHeader(http.StatusNoContent) // 204 No Content response
}

// Register User (Temporary JSON Storage)
func (apiCfg *ApiConfig) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into DB
	newUser, err := apiCfg.DB.RegisterUser(ctx, database.RegisterUserParams{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Age:          sql.NullInt32{Int32: req.Age, Valid: true},
	})
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("User created succesfully!")

	// Return the created user (excluding password)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

// Login User
func (apiCfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	user, err := apiCfg.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password
	if err := CheckPasswordHash(user.PasswordHash, req.Password); err != nil {
		log.Printf("%s", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	accessToken, err := apiCfg.MakeJWT(user.ID, time.Hour)
	log.Printf("Generated Token: %s", accessToken)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating refesh token: %s", err)
		w.WriteHeader(500)
		return
	}

	// Expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	params := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	_, err = apiCfg.DB.CreateRefreshToken(context.Background(), params)
	if err != nil {
		log.Printf("Error creating refesh token: %s", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("User login succesful!")

	req = User{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	returnData, err := json.Marshal(req)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Write(returnData)
}

func (apiCfg *ApiConfig) LogoutUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := apiCfg.DB.DeleteSession(context.Background(), userID)
	if err != nil {
		log.Printf("Error deleting refresh token: %s", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("User logged out succesfully, refresh token deleted!")
	w.WriteHeader(204)
}

// Edit User (Check JSON Data)
func (apiCfg *ApiConfig) EditUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("%s", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Handle password hashing only if provided
	var passwordHash string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		hashedPasswordStr := string(hashedPassword)
		passwordHash = hashedPasswordStr
	}

	// Update user
	updatedUser, err := apiCfg.DB.EditUser(ctx, database.EditUserParams{
		ID:           userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Age:          sql.NullInt32{Int32: req.Age, Valid: true},
		PasswordHash: passwordHash,
	})
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("User account updated succesfully!")

	// Return updated user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// Delete User (Check JSON Data)
func (apiCfg *ApiConfig) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := apiCfg.DB.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("%s", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	log.Printf("User deleted succesfully!")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (apiConfig *ApiConfig) PostRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(401)
		return
	}

	err = apiConfig.DB.RevokeRefreshToken(context.Background(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		w.WriteHeader(500)
	}

	log.Printf("Refresh Token revoked succesfully!")
	w.WriteHeader(204)
}

func (apiCfg *ApiConfig) PostRefresh(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the request context (already validated in middleware)
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	newAccessToken, err := apiCfg.MakeJWT(userID, time.Hour)
	log.Printf("Generated Token: %s", newAccessToken)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := MakeRefreshToken()
	if err != nil {
		log.Printf("Error generating refesh token: %s", err)
		w.WriteHeader(500)
		return
	}

	// Expires in 30 days
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	params := database.CreateRefreshTokenParams{
		Token:     newRefreshToken,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	_, err = apiCfg.DB.CreateRefreshToken(context.Background(), params)
	if err != nil {
		log.Printf("Error creating refesh token: %s", err)
		w.WriteHeader(500)
		return
	}

	user := User{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	returnData, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("Refresh token refreshed!")
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(returnData)
}

func (apiCfg *ApiConfig) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.DB.GetUsers(context.Background())
	if err != nil {
		log.Printf("Error retrieving users: %s", err)
		w.WriteHeader(500)
		return
	}

	var userList []User
	for _, user := range users {
		userList = append(userList, User{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Age:       user.Age.Int32,
		})
	}

	returnData, err := json.Marshal(userList)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, "Failed to process user data", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(returnData)
}

// Start Workout (Dummy Logic)
func StartWorkout(w http.ResponseWriter, r *http.Request) {
}

// Stop Workout (Dummy Logic)
func StopWorkout(w http.ResponseWriter, r *http.Request) {
}

// Get Workout Logs (Dummy Data)
func GetWorkoutLogs(w http.ResponseWriter, r *http.Request) {}
