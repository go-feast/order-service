

CREATE TYPE order_state as  ENUM (
    'order.canceled',
    'order.created',
    'order.paid',
    'order.cooking',
    'order.cooking.finished',
    'order.waiting',
    'order.taken',
    'order.delivering',
    'order.delivered',
    'order.closed'
);

CREATE TABLE orders(
                       id uuid PRIMARY KEY,
                       restaurant_id uuid,
                       customer_id uuid,
                       courier_id uuid,
                       meals uuid[],
                       state order_state,
                       destination geography(POINT),
                       transaction_id uuid NULL,
                       created_at timestamp
);

-- Add a comment to describe the 'orders' table
COMMENT ON TABLE orders IS 'This table represents an order. The state of the order is represented by the state column.';

-- Add a comment to describe the 'restaurant_id' column
-- Maybe change to a table (with subscription on restaurant.created)
COMMENT ON COLUMN orders.restaurant_id IS 'Unique identifier for the restaurant associated with the order.';

-- Add a comment to describe the 'customer_id' column
-- Maybe change to a table (with subscription on customer.created)
COMMENT ON COLUMN orders.customer_id IS 'Unique identifier for the customer who placed the order.';

-- Add a comment to describe the 'courier_id' column
COMMENT ON COLUMN orders.courier_id IS 'Unique identifier for the courier assigned to deliver the order.';

-- Add a comment to describe the 'meals' column
COMMENT ON COLUMN orders.meals IS 'Array of unique identifiers for the meals selected by the customer in the order.';

-- Add a comment to describe the 'state' column
COMMENT ON COLUMN orders.state IS 'Current state of the order in the order lifecycle, defined by the custom type order_state.';

-- Add a comment to describe the 'destination' column
COMMENT ON COLUMN orders.destination IS 'Geographic point (longitude/latitude) representing the delivery location for the order.';

-- Add a comment to describe the 'transaction_id' column
COMMENT ON COLUMN orders.transaction_id IS 'Unique identifier for the payment transaction associated with the order. Can be null if the transaction is not completed or not applicable.';

-- Add a comment to describe the 'created_at' column
COMMENT ON COLUMN orders.created_at IS 'Timestamp indicating when the order was created.';
