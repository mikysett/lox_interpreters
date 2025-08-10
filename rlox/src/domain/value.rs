pub enum Value {
    Double(f64),
}

impl Value {
    pub fn print(&self) {
        match self {
            Value::Double(n) => print!("{n}"),
        }
    }
}
