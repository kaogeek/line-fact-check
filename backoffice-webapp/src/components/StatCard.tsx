import { Card, CardContent } from './ui/card';
import { TYH2, TYP } from './Typography';

interface StatCardProp {
  label: string;
  value: string;
}

export default function StatCard({ label, value }: StatCardProp) {
  return (
    <Card>
      <CardContent>
        <TYH2>{value}</TYH2>
        <TYP>{label}</TYP>
      </CardContent>
    </Card>
  );
}
