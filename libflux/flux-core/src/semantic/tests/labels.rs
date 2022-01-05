use super::*;

#[test]
fn labels() {
    test_infer! {
        env: map![
            "fill" => "(<-tables: [{ A with D: B }], ?column: D, ?value: B, ?usePrevious: bool) => [C]
                where C: Record,
                      D: Label
                "
        ],
        src: r#"
            x = [{ a: 1 }] |> fill(column: "a", value: 2)
        "#,
        exp: map![
            "x" => "{ a: int }",
        ],
    }
}
