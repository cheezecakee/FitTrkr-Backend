## Routes

| Method | Pattern                              | Handler                    | Action |
| ------ | ------------------------------------ | -------------------------- | ------ |
| POST   | /api/user/register                   | RegisterUser               |        |
| POST   | /api/user/login                      | LoginUser                  |        |
| POST   | /api/user/logout                     | LogoutUser                 |        |
| POST   | /api/user/revoke                     | PostRevoke                 |        |
| POST   | /api/user/refresh                    | PostRefresh                |        |
| PUT    | /api/user/edit                       | EditUser                   |        |
| DELETE | /api/user                            | DeleteUser                 |        |
| GET    | /api/workouts                        | GetWorkouts                |        |
| POST   | /api/workouts                        | CreateWorkout              |        |
| PUT    | /api/workouts/{id}                   | EditWorkout                |        |
| DELETE | /api/workouts/{id}                   | DeleteWorkout              |        |
| POST   | /api/workouts/exercises              | CreateWorkoutExercise      |        |
| PUT    | /api/workouts/exercises/{id}         | EditWorkoutExercise        |        |
| DELETE | /api/workouts/exercises/{id}         | DeleteWorkoutExercise      |        |
| POST   | /api/session/start                   | SetWorkoutSession          |        |
| GET    | /api/session                         | GetActiveWorkoutSession    |        |
| PUT    | /api/session/stop                    | StopActiveWorkoutSession   |        |
| PUT    | /api/session/finish                  | FinishActiveWorkoutSession |        |
| GET    | /api/session/workout                 | GetWorkoutSessionDetails   |        |
| PUT    | /api/session/workout/finish          | FinishWorkoutSession       |        |
| POST   | /api/session/workout/exercise/start  | StartExercise              |        |
| PUT    | /api/session/workout/exercise/log    | LogExerciseSet             |        |
| PUT    | /api/session/workout/exercise/stop   | StopExercise               |        |
| PUT    | /api/session/workout/exercise/finish | FinishExercise             |        |
|        |                                      |                            |        |
|        |                                      |                            |        |
|        |                                      |                            |        |
|        |                                      |                            |        |
