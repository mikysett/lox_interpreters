use crate::scanner::Token;
use std::fmt;

#[derive(Debug, Clone)]
pub struct MetaObject {
    pub obj: Object,
}

#[derive(Debug, Clone)]
pub enum Object {
    Str(String),
}

impl fmt::Display for MetaObject {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match &self.obj {
            Object::Str(n) => write!(f, "{n}"),
        }
    }
}
pub fn copy_string(previous: &Token) -> MetaObject {
    let full_str = previous.as_str();
    MetaObject {
        obj: Object::Str(full_str[1..full_str.len() - 1].to_owned()),
    }
}

// `is_string` is to be used in matching arms to verify if a `Value` is an object of type `Str`
#[macro_export]
macro_rules! is_string {
    ($name:ident) => {
        Value::Obj(MetaObject {
            obj: Object::Str($name),
        })
    };
}
