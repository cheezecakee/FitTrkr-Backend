## Routes
### User 
POST   /api/user/register
POST   /api/user/login	
POST   /api/user/logout
POST   /api/user/revoke
POST   /api/user/refresh
PUT    /api/user/edit"
DELETE /api/user"

### Workout 
GET    /api/workouts	
POST   /api/workouts
PUT    /api/workouts/{id}

### WorkoutExercise 
POST   /api/workouts/exercises
PUT    /api/workouts/exercises/{id}
DELETE /api/workouts/exercises/{id}
    
### Exercises
GET    /api/exercises
GET    /api/exercises/{id}

### Start/Stop Workout
POST   /api/workouts/start
POST   /api/workouts/stop

### Logs 
GET    /api/logs


