module backend/task-tracker

go 1.25.3

replace backend/jsondatabase => ./json-database

replace backend/tasks => ./tasks

require (
	backend/jsondatabase v0.0.0-00010101000000-000000000000 // indirect
	backend/tasks v0.0.0-00010101000000-000000000000
)
