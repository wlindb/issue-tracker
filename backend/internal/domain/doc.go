// Package domain contains the core types, interfaces, and domain services for the issue tracker.
// Domain entities (user.go) have zero external dependencies.
// Domain services (service.go) may import external libraries (e.g. bcrypt, JWT) for business logic,
// but must not import any other internal/ package.
package domain
