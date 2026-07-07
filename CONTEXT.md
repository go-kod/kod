# Kod

Kod is a dependency injection context for Go applications. Its language distinguishes dependency wiring from application concerns that should remain in user code.

## Language

**Component**:
A user-defined interface and implementation pair managed by Kod.
_Avoid_: Service, provider

**Business Configuration**:
Application settings owned and loaded by user components, not by Kod core.
_Avoid_: Kod config, framework config
