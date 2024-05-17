import { connect } from "nats";
import fs from "fs/promises";
import crypto from "crypto";

async function main() {
    const server = { servers: "localhost:4222" };
    const nc = await connect(server);
    const js = nc.jetstream();
    const data = await fs.readFile("model.json");
    for (let i = 0; i < 10000; i++) {
        let order = JSON.parse(data);
        order["order_uid"] = crypto
            .createHash("sha256")
            .update(i.toString())
            .digest("hex");
        const newData = JSON.stringify(order);
        let pa = await js.publish("ORDERS.new", newData);
        console.log(pa, order["order_uid"]);
    }
    return;
}
main();
