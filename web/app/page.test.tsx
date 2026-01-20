import { render, screen, waitFor } from '@testing-library/react'
import Dashboard from './page'
import { getRecentRuns } from '@/lib/api'

// Mock the API module
jest.mock('@/lib/api', () => ({
    getRecentRuns: jest.fn(),
}))

describe('Dashboard', () => {
    it('renders loading state initially', () => {
        (getRecentRuns as jest.Mock).mockReturnValue(new Promise(() => { })) // Never resolves
        render(<Dashboard />)
        expect(screen.getByText('Loading...')).toBeInTheDocument()
    })

    it('renders data after loading', async () => {
        const mockData = [
            {
                source_id: 'test_src',
                status: 'FAIL',
                records_checked: 100,
                rules_failed: 5,
                timestamp: new Date().toISOString(),
            },
        ]
            ; (getRecentRuns as jest.Mock).mockResolvedValue(mockData)

        render(<Dashboard />)

        await waitFor(() => {
            expect(screen.getByText('DataGuard Dashboard')).toBeInTheDocument()
            expect(screen.getByText('test_src')).toBeInTheDocument()
            expect(screen.getByText('FAIL')).toBeInTheDocument()
        })
    })

    it('renders empty state if no data', async () => {
        (getRecentRuns as jest.Mock).mockResolvedValue([])

        render(<Dashboard />)

        await waitFor(() => {
            expect(screen.getByText('No matching records found')).toBeInTheDocument()
        })
    })
})
