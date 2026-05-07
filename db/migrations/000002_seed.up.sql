-- Inserindo Quarto
INSERT INTO rooms (id, number, type, capacity, price_per_night, status) 
VALUES ('r1-uuid-placeholder', '101', 'STANDARD', 2, 150.00, 'ATIVO');

-- Inserindo Hóspede
INSERT INTO guests (id, full_name, document, email, phone) 
VALUES ('g1-uuid-placeholder', 'João Silva', '12345678900', 'joao@example.com', '11999999999');
