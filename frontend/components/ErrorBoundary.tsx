'use client';

import { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    // Suppress router mounting errors (known Next.js dev mode issue with HotReload)
    const errorMsg = error?.message || '';
    if (errorMsg.includes('expected app router to be mounted') || 
        errorMsg.includes('invariant expected app router')) {
      // This is a Next.js internal error during hot reload, ignore it
      return { hasError: false };
    }
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: any) {
    // Suppress router mounting errors in console (Next.js HotReload issue)
    const errorMsg = error?.message || '';
    if (errorMsg.includes('expected app router to be mounted') || 
        errorMsg.includes('invariant expected app router')) {
      // Silently ignore - this is a Next.js dev mode hot reload issue
      return;
    }
    console.error('ErrorBoundary caught an error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError && this.state.error) {
      return this.props.fallback || <div>Something went wrong.</div>;
    }

    return this.props.children;
  }
}
