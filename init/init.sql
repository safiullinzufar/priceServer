CREATE TABLE pricetable (
    mail TEXT,
    link TEXT,
    price VARCHAR(32),
    unique(mail, link)
);