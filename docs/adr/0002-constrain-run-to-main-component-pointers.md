# Constrain Run to Main Component Pointers

Kod's `Run` and `MustRun` APIs accept the application entry point through a generic run function. We constrain that generic type to a main component implementation pointer so invalid value-type entry points fail at compile time instead of reaching `Get` and failing at runtime, while keeping normal `kod.Run(ctx, func(ctx context.Context, app *App) error)` calls inferable.
