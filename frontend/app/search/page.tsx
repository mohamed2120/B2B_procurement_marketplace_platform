'use client';

import { useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { apiClients } from '@/lib/api';
import { isAuthenticated } from '@/lib/auth';
import Card from '@/components/ui/Card';
import Button from '@/components/ui/Button';
import Link from 'next/link';

interface SearchResult {
  type: string;
  id: string;
  title: string;
  description?: string;
  fields: Record<string, any>;
  score?: number;
}

interface SearchResponse {
  results: SearchResult[];
  facets: Record<string, any>;
  total: number;
  page: number;
  page_size: number;
}

export default function SearchPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const query = searchParams.get('q') || '';
  const typeParam = searchParams.get('type') || 'all';
  const pageParam = parseInt(searchParams.get('page') || '1');
  const sortParam = searchParams.get('sort') || 'relevance';

  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(pageParam);
  const [pageSize, setPageSize] = useState(10);
  const [searchType, setSearchType] = useState(typeParam);
  const [sort, setSort] = useState(sortParam);
  const [isGuest, setIsGuest] = useState(true);
  const [searchQuery, setSearchQuery] = useState(query);

  useEffect(() => {
    setIsGuest(!isAuthenticated());
  }, []);

  useEffect(() => {
    if (query) {
      performSearch();
    }
  }, [query, searchType, page, sort]);

  const performSearch = async () => {
    if (!query.trim()) {
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await apiClients.search.get<SearchResponse>('/api/v1/search', {
        params: {
          q: query,
          type: searchType,
          page,
          page_size: isGuest ? 10 : 20,
          sort,
        },
      });

      setResults(response.data.results || []);
      setTotal(response.data.total || 0);
      setPageSize(response.data.page_size || 10);
    } catch (err: any) {
      if (err.response?.status === 503 || err.message?.includes('opensearch')) {
        setError('Search is temporarily unavailable. Please try again later.');
      } else {
        setError(err.response?.data?.error || 'Search failed. Please try again.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleTypeChange = (newType: string) => {
    setSearchType(newType);
    setPage(1);
    router.push(`/search?q=${encodeURIComponent(query)}&type=${newType}&sort=${sort}`);
  };

  const handleSortChange = (newSort: string) => {
    setSort(newSort);
    setPage(1);
    router.push(`/search?q=${encodeURIComponent(query)}&type=${searchType}&sort=${newSort}`);
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
    router.push(`/search?q=${encodeURIComponent(query)}&type=${searchType}&sort=${sort}&page=${newPage}`);
  };

  const getResultLink = (result: SearchResult) => {
    switch (result.type) {
      case 'part':
        return `/app/catalog/parts/${result.id}`;
      case 'equipment':
        return `/app/equipment/${result.id}`;
      case 'company':
        return `/app/companies/${result.id}`;
      case 'listing':
        return `/app/marketplace/listings/${result.id}`;
      default:
        return '#';
    }
  };

  const renderResult = (result: SearchResult) => {
    const typeLabels: Record<string, string> = {
      part: 'Part',
      equipment: 'Equipment',
      company: 'Company',
      listing: 'Product',
      service: 'Service',
    };

    return (
      <Card key={`${result.type}-${result.id}`} className="mb-4 hover:shadow-md transition-shadow">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <div className="flex items-center space-x-2 mb-2">
              <span className="px-2 py-1 text-xs font-semibold rounded bg-primary-100 text-primary-800">
                {typeLabels[result.type] || result.type}
              </span>
              {result.score && (
                <span className="text-xs text-gray-500">Relevance: {result.score.toFixed(2)}</span>
              )}
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-1">
              <Link href={getResultLink(result)} className="hover:text-primary-600">
                {result.title}
              </Link>
            </h3>
            {result.description && (
              <p className="text-gray-600 text-sm mb-2">{result.description}</p>
            )}
            <div className="flex flex-wrap gap-2 text-sm text-gray-500">
              {result.fields.manufacturer && (
                <span>Manufacturer: {result.fields.manufacturer}</span>
              )}
              {result.fields.category && (
                <span>Category: {result.fields.category}</span>
              )}
              {result.fields.price && !result.fields.price_restricted && (
                <span>Price: {result.fields.currency || '$'} {result.fields.price}</span>
              )}
              {result.fields.price_restricted && (
                <span className="text-primary-600">Contact for pricing</span>
              )}
            </div>
          </div>
          <Link href={getResultLink(result)}>
            <Button size="sm" variant="secondary">View Details</Button>
          </Link>
        </div>
        {isGuest && result.type === 'listing' && (
          <div className="mt-3 pt-3 border-t border-gray-200">
            <p className="text-sm text-gray-600">
              <Link href="/login" className="text-primary-600 hover:underline">
                Login
              </Link>
              {' '}to see full details, stock availability, and pricing
            </p>
          </div>
        )}
      </Card>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Search Results</h1>

        {isGuest && (
          <div className="mb-6 bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <p className="text-sm text-yellow-800">
              <strong>Limited Results:</strong> You're viewing public results only. 
              <Link href="/login" className="text-yellow-900 underline ml-1">
                Login
              </Link>
              {' '}to see more results and full details.
            </p>
          </div>
        )}

        {error && (
          <div className="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}

        {query && (
          <>
            {/* Type Tabs */}
            <div className="mb-6 flex space-x-1 border-b border-gray-200">
              {['all', 'part', 'equipment', 'company', 'listing'].map((t) => (
                <button
                  key={t}
                  onClick={() => handleTypeChange(t)}
                  className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
                    searchType === t
                      ? 'border-primary-500 text-primary-600'
                      : 'border-transparent text-gray-500 hover:text-gray-700'
                  }`}
                >
                  {t.charAt(0).toUpperCase() + t.slice(1)}
                </button>
              ))}
            </div>

            {/* Sort and Filters */}
            <div className="mb-6 flex justify-between items-center">
              <div className="text-sm text-gray-600">
                {total > 0 ? (
                  <>Found {total} result{total !== 1 ? 's' : ''} for "{query}"</>
                ) : (
                  <>No results found for "{query}"</>
                )}
              </div>
              <select
                value={sort}
                onChange={(e) => handleSortChange(e.target.value)}
                className="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                <option value="relevance">Relevance</option>
                <option value="rating">Rating</option>
                <option value="price">Price (Low to High)</option>
                <option value="eta">ETA (Fastest)</option>
              </select>
            </div>

            {/* Results */}
            {loading ? (
              <div className="text-center py-12">
                <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto mb-4"></div>
                <p className="text-gray-600">Searching...</p>
              </div>
            ) : results.length > 0 ? (
              <>
                {results.map(renderResult)}
                
                {/* Pagination */}
                {total > pageSize && (
                  <div className="mt-8 flex justify-center space-x-2">
                    <Button
                      variant="secondary"
                      onClick={() => handlePageChange(page - 1)}
                      disabled={page === 1}
                    >
                      Previous
                    </Button>
                    <span className="px-4 py-2 text-sm text-gray-600">
                      Page {page} of {Math.ceil(total / pageSize)}
                    </span>
                    <Button
                      variant="secondary"
                      onClick={() => handlePageChange(page + 1)}
                      disabled={page >= Math.ceil(total / pageSize)}
                    >
                      Next
                    </Button>
                  </div>
                )}
              </>
            ) : !loading && query ? (
              <Card>
                <p className="text-center text-gray-600 py-8">
                  No results found. Try different keywords or{' '}
                  <Link href="/login" className="text-primary-600 hover:underline">
                    login
                  </Link>
                  {' '}to see more results.
                </p>
              </Card>
            ) : null}
          </>
        )}

        {!query && (
          <Card>
            <p className="text-center text-gray-600 py-8">
              Enter a search query to find parts, equipment, companies, and listings.
            </p>
          </Card>
        )}
      </div>
    </div>
  );
}
