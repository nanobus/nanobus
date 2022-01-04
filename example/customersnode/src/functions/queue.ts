export interface Queue<T = any> {
  push(task: T): void;
  receive(): Promise<T | undefined>;
  close(): void;
}

export class queue<T = any> implements Queue {
  private q: T[] = [];
  private p: Promise<T> = undefined;
  private resolve: (value: T | PromiseLike<T>) => void = undefined;
  private reject: (reason?: any) => void = undefined;
  private closed: boolean = false;

  push(task: T): void {
    if (this.closed) {
      throw new Error("queue is closed");
    }
    const l = this.q.length;
    this.q.push();
    if (l === 0 && this.resolve !== undefined) {
      this.resolve(task);
      this.p = undefined;
      this.resolve = undefined;
      this.reject = undefined;
    }
  }

  close(): void {
    if (this.closed) {
      return;
    }
    if (this.resolve !== undefined) {
      this.resolve(undefined);
    }
    this.closed = true;
  }

  async receive(): Promise<T | undefined> {
    if (this.q.length > 0) {
      return this.q.shift();
    }
    if (this.closed) {
      return undefined;
    }
    const that = this;
    this.p = new Promise((resolve, reject) => {
      that.resolve = resolve;
      that.reject = reject;
    });
    return this.p;
  }
}
