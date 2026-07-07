# Keep Configuration Out of Core

Kod's core is dependency injection: component wiring, references, lifecycle, testing fakes, and interceptors. We removed built-in config loading so application configuration is modeled as ordinary business components, avoiding a permanent dependency on one config stack and keeping Kod closer to Wire/Fx-style DI boundaries.
