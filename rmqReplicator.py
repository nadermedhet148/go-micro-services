import pika
import sys

def consume_and_produce(src_rabbitmq_url, dest_rabbitmq_url, queue_name):
    # Connect to source RabbitMQ
    src_params = pika.URLParameters(src_rabbitmq_url)
    src_connection = pika.BlockingConnection(src_params)
    src_channel = src_connection.channel()

    dest_params = pika.URLParameters(dest_rabbitmq_url)
    dest_connection = pika.BlockingConnection(dest_params)
    dest_channel = dest_connection.channel()

    # id need to declare the queue on the destination RabbitMQ
    # dest_channel.queue_declare(queue=queue_name, durable=True)

    def callback(ch, method, properties, body):
        print(f"Consuming message from {src_rabbitmq_url}: {body}")
        dest_channel.basic_publish(
            exchange='',
            routing_key=queue_name,
            body=body,
            properties=pika.BasicProperties(delivery_mode=2)
        )
        ch.basic_ack(delivery_tag=method.delivery_tag)

    src_channel.basic_consume(queue=queue_name, on_message_callback=callback)

    print(f"Waiting for messages in {queue_name}. To exit press CTRL+C")
    try:
        src_channel.start_consuming()
    except KeyboardInterrupt:
        print("Stopping...")
    finally:
        src_connection.close()
        dest_connection.close()

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python script.py <src_rabbitmq_url> <dest_rabbitmq_url> <queue_name>")
        sys.exit(1)

    src_rabbitmq_url = sys.argv[1]
    dest_rabbitmq_url = sys.argv[2]
    queue_name = sys.argv[3]

    consume_and_produce(src_rabbitmq_url, dest_rabbitmq_url, queue_name)
