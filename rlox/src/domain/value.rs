use std::fmt;
#[derive(Debug, Clone)]
#[repr(u8)] // TODO: Verify this doesn't break with values > u8
pub enum Value {
    Unknown = 0,
    Double(f64),
}

impl fmt::Display for Value {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Value::Double(n) => write!(f, "{n}"),
            Value::Unknown => write!(f, "unknown"),
        }
    }
}
