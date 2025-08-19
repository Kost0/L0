CREATE TABLE delivery (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50) NOT NULL,
    zip VARCHAR(50) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    order_uid UUID PRIMARY KEY,
    track_number VARCHAR(255) NOT NULL UNIQUE,
    entry VARCHAR(50) NOT NULL,
    delivery_id UUID REFERENCES delivery(id),
    locale VARCHAR(10) NOT NULL,
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255) NOT NULL,
    delivery_service VARCHAR(255) NOT NULL,
    shardkey VARCHAR(255) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard VARCHAR(255) NOT NULL
);

CREATE TABLE payment (
     transaction UUID REFERENCES orders(order_uid),
     request_id VARCHAR(255),
     currency VARCHAR(3) NOT NULL,
     provider VARCHAR(255) NOT NULL,
     amount INT NOT NULL,
     payment_dt BIGINT NOT NULL,
     bank VARCHAR(255) NOT NULL,
     delivery_cost INT NOT NULL,
     goods_total INT NOT NULL,
     custom_fee INT NOT NULL
);

CREATE TABLE items (
    chrt_id INT NOT NULL,
    track_number VARCHAR(255) REFERENCES orders(track_number),
    price INT NOT NULL,
    rid VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(50) NOT NULL,
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INT NOT NULL
);