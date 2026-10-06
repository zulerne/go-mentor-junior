 create table orders (
    id uuid primary key,
    customer_id uuid not null,
    restaurant_id uuid not null,
    status text not null check (
        status in (
            'PENDING',
            'ACCEPTED',
            'PREPARING',
            'READY_FOR_PICKUP',
            'REJECTED',
            'CANCELLED'
        )
    ),
    subtotal_minor bigint not null check (subtotal_minor >= 0),
    currency text not null,
    delivery_address text not null,
    rejection_reason text not null,
    delivery_status text check (
        delivery_status in (
            'WAITING_FOR_PREPARATION',
            'IN_DELIVERY'
        )
    ),
    created_at timestamptz not null,
    updated_at timestamptz not null
);

create index orders_customer_id_idx on orders (customer_id);
create index orders_restaurant_id_idx on orders (restaurant_id);

create table order_items (
    id uuid primary key,
    order_id uuid not null references orders(id),
    menu_item_id uuid not null,
    name text not null,
    unit_price_minor bigint not null check (unit_price_minor >= 0),
    quantity int not null check (quantity between 1 and 10),
    instructions text not null default '',

    unique (order_id, menu_item_id)
);

create table carts (
    customer_id uuid primary key,
    restaurant_id uuid,
    currency text,
    subtotal_minor bigint not null check (subtotal_minor >= 0) default 0
);

create table cart_items (
    id uuid primary key,
    customer_id uuid not null references carts(customer_id),
    menu_item_id uuid not null,
    name text not null,
    unit_price_minor bigint not null check (unit_price_minor >= 0),
    currency text not null,
    quantity int not null check (quantity between 1 and 10),
    instructions text not null default '',

    unique (customer_id, menu_item_id)
);
