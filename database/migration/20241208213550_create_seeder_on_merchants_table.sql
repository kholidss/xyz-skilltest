-- +goose Up
-- +goose StatementBegin
INSERT INTO merchants (id, name, slug)
VALUES
    ('1b4e28ba-2fa1-4b21-a14a-64533fe21bc1', 'PT XYZ', 'pt-xyz'),
    ('2c4e28ba-3fa2-4b31-b15b-74533fe22cd2', 'Asahi Shop', 'asahi-sop'),
    ('5686a5b6-0ad6-4588-9843-abb0b9b47bd6', 'Djaya Dealer', 'djaya-dealer');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM merchants
WHERE id IN ('1b4e28ba-2fa1-4b21-a14a-64533fe21bc1', '2c4e28ba-3fa2-4b31-b15b-74533fe22cd2', '5686a5b6-0ad6-4588-9843-abb0b9b47bd6');
-- +goose StatementEnd
