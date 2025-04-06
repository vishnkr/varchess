/*import init, { ChessCoreLib } from 'stonkfish';

class ChessCoreWrapper {
    private coreInstance: ChessCoreLib | null = null;
  
    async initWasm() {
      await init();
    }
  
    loadPosition(config_json: string) {
      this.coreInstance = new ChessCoreLib(config_json);
    }
  
    getLegalMoves() {
      return this.coreInstance ? this.coreInstance.getLegalMoves() : [];
    }
  
    makeMove(move: string) {
      if (this.coreInstance) {
        //this.coreInstance.makeMove(move);
      }
    }
  
    // Add other methods as needed
  }
  
  const chessCoreWrapper = new ChessCoreWrapper();
  export default chessCoreWrapper;
  		"wasm": "wasm-pack build ./stonkfish --target web"
  */