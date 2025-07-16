import { useTranslation } from 'react-i18next';
import { Card, CardContent, CardHeader, CardTitle } from '../../components/ui/card';
import StatCard from '../../components/StatCard';
import { TYH3 } from '../../components/Typography';
import {
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

// Sample data for charts
const statusData = [
  { name: 'Jan', value: 2 },
  { name: 'Feb', value: 3 },
  { name: 'Mar', value: 5 },
  { name: 'Apr', value: 4 },
  { name: 'May', value: 7 },
  { name: 'Jun', value: 9 },
];

const approveData = [
  { name: 'Real', value: 30 },
  { name: 'Fake', value: 70 },
];

const COLORS = ['#0088FE', '#FF8042'];

export default function DashboardPage() {
  const { t } = useTranslation();

  return (
    <div className="p-4 h-full space-y-4">
      <div>
        <TYH3 className="mb-4">{t('dashboard.title')}</TYH3>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <StatCard label={t('dashboard.totalTicket')} value="9" />
        <StatCard label={t('dashboard.pending')} value="3" />
        <StatCard label={t('dashboard.answer')} value="2" />
        <StatCard label={t('dashboard.rejected')} value="2" />
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>{t('dashboard.statusTimeline')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[300px]">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={statusData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Line type="monotone" dataKey="value" stroke="#8884d8" activeDot={{ r: 8 }} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>{t('dashboard.approveStat')}</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-3 gap-4 mb-4">
              <StatCard label={t('dashboard.totalNews')} value="9" />
              <StatCard label={t('dashboard.realNews')} value="9" />
              <StatCard label={t('dashboard.fakeNews')} value="3" />
            </div>
            <div className="h-[200px]">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={approveData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="value"
                    label={({ name, percent }: { name: string; percent?: number }) =>
                      `${name} ${percent ? (percent * 100).toFixed(0) : 0}%`
                    }
                  >
                    {approveData.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
