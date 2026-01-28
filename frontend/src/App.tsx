const App = () => {
  return (
    <div className="app">
      <header className="app__header">
        <h1>G1STI Windows Client</h1>
        <p>Minimal client for VLESS subscriptions powered by sing-box.</p>
      </header>
      <section className="app__card">
        <h2>Disclaimer</h2>
        <p>
          This software is intended strictly for lawful, authorized network
          administration and corporate use. You are responsible for complying
          with local laws, policies, and provider terms. We do not endorse or
          support any misuse.
        </p>
      </section>
      <section className="app__card">
        <h2>Status</h2>
        <p>Backend scaffolding in progress. UI will be wired in Step 3.</p>
      </section>
    </div>
  );
};

export default App;
