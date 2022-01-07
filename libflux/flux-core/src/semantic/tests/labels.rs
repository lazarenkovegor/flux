use super::*;

#[test]
fn labels() {
    test_infer! {
        env: map![
            "fill" => "(<-tables: [{ A with B: C }], ?column: B, ?value: D) => [{ A with B: D }]
                where B: Label
                "
        ],
        src: r#"
            x = [{ a: 1 }] |> fill(column: "a", value: "x")
            y = [{ a: 1, b: ""}] |> fill(column: "b", value: 1.0)
        "#,
        exp: map![
            "x" => "[{ a: string }]",
            "y" => "[{ a: int, b: float }]",
        ],
    }
}
