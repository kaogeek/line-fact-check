import { TYH3 } from './Typography';

interface SectionTopicProps {
  label: string;
}

export default function SectionTopic({ label }: SectionTopicProps) {
  return (
    <div className="flex flex-col p-4">
      <TYH3>{label}</TYH3>
    </div>
  );
}
