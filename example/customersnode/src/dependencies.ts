import logger from './logger';
import dotenv from 'dotenv';
import { HTTPHandlers, HTTPInvoker, msgpackCodec } from './functions';
import { OutboundImpl, Service } from './interfaces';

const result = dotenv.config();
if (result.error) {
  dotenv.config({ path: '.env.default' });
}

const PORT = parseInt(process.env.PORT) || 3000;
const HOST = process.env.HOST;

const myInvoker = HTTPInvoker(
  process.env.OUTBOUND_URL || 'http://localhost:32321/outbound',
  msgpackCodec
);

export const outbound = new OutboundImpl(myInvoker);

export const handlers = new HTTPHandlers();
export const service = new Service(handlers, msgpackCodec, () => {
  handlers.listen(PORT, HOST, () => {
    logger.info(`🌏 Nanoprocess server started at http://${HOST}:${PORT}`);
  });
});