import { SearchParams } from './search-params';

describe('SearchParams Unit Tests', () => {
  describe('page prop', () => {
    const testCases = [
      { page: null, expected: 1 },
      { page: undefined, expected: 1 },
      { page: '', expected: 1 },
      { page: 'fake', expected: 1 },
      { page: 0, expected: 1 },
      { page: -1, expected: 1 },
      { page: 5.5, expected: 1 },
      { page: true, expected: 1 },
      { page: false, expected: 1 },
      { page: {}, expected: 1 },

      { page: 1, expected: 1 },
      { page: 2, expected: 2 },
    ];

    test.each(testCases)('returns correct output when page is = %j', ({ page, expected }) => {
      const params = new SearchParams();
      expect(params.page).toBe(1);
      expect(new SearchParams({ page } as any).page).toBe(expected);
    });
  });

  describe('per_page prop', () => {
    const testCases = [
      { per_page: null, expected: 15 },
      { per_page: undefined, expected: 15 },
      { per_page: '', expected: 15 },
      { per_page: 'fake', expected: 15 },
      { per_page: 0, expected: 15 },
      { per_page: -1, expected: 15 },
      { per_page: 5.5, expected: 15 },
      { per_page: true, expected: 15 },
      { per_page: false, expected: 15 },
      { per_page: {}, expected: 15 },

      { per_page: 1, expected: 1 },
      { per_page: 2, expected: 2 },
      { per_page: 10, expected: 10 },
    ];

    test.each(testCases)(
      'returns correct output when per_page is = %j',
      ({ per_page, expected }) => {
        const params = new SearchParams();
        expect(params.per_page).toBe(15);
        expect(new SearchParams({ per_page } as any).per_page).toBe(expected);
      },
    );
  });

  describe('sort prop', () => {
    const testCase = [
      { sort: null, expected: null },
      { sort: undefined, expected: null },
      { sort: '', expected: null },
      { sort: 0, expected: '0' },
      { sort: -1, expected: '-1' },
      { sort: 5.5, expected: '5.5' },
      { sort: true, expected: 'true' },
      { sort: false, expected: 'false' },
      { sort: {}, expected: '[object Object]' },
      { sort: 'field', expected: 'field' },
    ];

    test.each(testCase)('returns correct output when sort is = %j', ({ sort, expected }) => {
      const params = new SearchParams();
      expect(params.sort).toBeNull();
      expect(new SearchParams({ sort } as any).sort).toBe(expected);
    });
  });

  describe('sort_dir prop', () => {
    const testCases = [
      { sort_dir: null, expected: 'asc' },
      { sort_dir: undefined, expected: 'asc' },
      { sort_dir: '', expected: 'asc' },
      { sort_dir: 0, expected: 'asc' },
      { sort_dir: 'fake', expected: 'asc' },

      { sort_dir: 'asc', expected: 'asc' },
      { sort_dir: 'ASC', expected: 'asc' },
      { sort_dir: 'desc', expected: 'desc' },
      { sort_dir: 'DESC', expected: 'desc' },
    ];

    test.each(testCases)(
      'returns correct output when {sort:"field"} and sort_dir is = %j',
      ({ sort_dir, expected }) => {
        expect(new SearchParams({ sort: 'field', sort_dir: sort_dir } as any).sort_dir).toBe(
          expected,
        );
      },
    );
  });

  describe('fields prop', () => {
    const testCases = [
      { filter: null, expected: null },
      { filter: undefined, expected: null },
      { filter: '', expected: null },

      { filter: 0, expected: '0' },
      { filter: -1, expected: '-1' },
      { filter: 5.5, expected: '5.5' },
      { filter: true, expected: 'true' },
      { filter: false, expected: 'false' },
      { filter: {}, expected: '[object Object]' },
      { filter: 'field', expected: 'field' },
    ];

    test.each(testCases)('returns correct output when fields is = %j', ({ filter, expected }) => {
      const params = new SearchParams();
      expect(params.filter).toBeNull();
      expect(new SearchParams({ filter: filter } as any).filter).toBe(expected);
    });
  });
});
