import { registerInbound, registerCustomerActor, outbound, start } from "./adapter";
import { CustomerActorImpl, InboundImpl } from "./service";

registerInbound(new InboundImpl(outbound));
registerCustomerActor(new CustomerActorImpl());

start();