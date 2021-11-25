import { outbound, registerCustomerActor, registerInbound, start } from "./adapter";
import { CustomerActorImpl, InboundImpl } from "./service";

registerInbound(new InboundImpl(outbound));
registerCustomerActor(new CustomerActorImpl());

start();