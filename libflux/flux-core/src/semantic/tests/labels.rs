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
        "#,
        exp: map![
            "x" => "[{ a: string }]",
        ],
    }
}
