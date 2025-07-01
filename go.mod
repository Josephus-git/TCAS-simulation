module github.com/josephus-git/TCAS-simulation

replace github.com/josephus-git/TCAS-simulation/internal/config => ./internal/config

replace github.com/josephus-git/TCAS-simulation/internal/aviation => ./internal/aviation

replace github.com/josephus-git/TCAS-simulation/internal/tcas => ./internal/tcas

require github.com/josephus-git/TCAS-simulation/internal/aviation v0.0.0

go 1.22.2
