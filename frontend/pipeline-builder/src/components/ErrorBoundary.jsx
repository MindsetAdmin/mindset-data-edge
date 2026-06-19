import { Component } from 'react';

// Keeps a render/effect error in one page from blanking the whole app.
export default class ErrorBoundary extends Component {
  constructor(props) {
    super(props);
    this.state = { error: null };
  }

  static getDerivedStateFromError(error) {
    return { error };
  }

  componentDidUpdate(prevProps) {
    // Reset the boundary when navigating to a different route.
    if (prevProps.resetKey !== this.props.resetKey && this.state.error) {
      this.setState({ error: null });
    }
  }

  render() {
    if (this.state.error) {
      return (
        <div className="h-full flex items-center justify-center p-6">
          <div className="bg-red-500/15 border border-red-500/40 rounded-lg p-5 max-w-md text-center">
            <div className="text-3xl mb-2">⚠️</div>
            <h2 className="text-white font-semibold mb-1">Une erreur est survenue</h2>
            <p className="text-red-300 text-sm break-words">{String(this.state.error.message || this.state.error)}</p>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}
