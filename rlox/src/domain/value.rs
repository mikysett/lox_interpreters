use std::fmt;
#[derive(Debug, Clone)]
#[repr(u8)] // TODO: Verify this doesn't break with values > u8
pub enum Value {
    Unknown = 0,
    Double(f64),
    Bool(bool),
    Nil,
}

impl Value {
    pub fn is_falsey(&self) -> bool {
        match self {
            Value::Nil => true,
            Value::Bool(false) => true,
            _ => false,
        }
    }

    pub fn equal(&self, other: &Value) -> bool {
        match (self, other) {
            (Value::Bool(a), Value::Bool(b)) => a == b,
            (Value::Nil, Value::Nil) => true,
            (Value::Double(a), Value::Double(b)) => a == b,
            _ => false,
        }
    }

    pub fn is_zero(&self) -> bool {
        match self {
            Value::Double(n) => *n == 0.0,
            _ => false,
        }
    }
}

impl fmt::Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Value::Double(n) => write!(f, "{n}"),
            Value::Bool(n) => write!(f, "{n}"),
            Value::Nil => write!(f, "nil"),
            Value::Unknown => write!(f, "unknown"),
        }
    }
}
