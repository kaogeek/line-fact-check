import { TYH3 } from '@/components/Typography';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Form, FormControl, FormField, FormItem } from '@/components/ui/form';
import { useTranslation } from 'react-i18next';
import { useEffect } from 'react';
import { useGetTopicAnswerByTopicId } from '@/hooks/api/topicAnswer';
import { TopicAnswerType } from '@/lib/api/type/topic-answer';
import LoadingState from '@/components/state/LoadingState';
import ErrorState from '@/components/state/ErrorState';
import { History } from 'lucide-react';

interface TopicMessageAnswerProps {
  onClickHistory: () => void;
  topicId: string;
  onUpdateAnswer: (answerId: string, content: string) => Promise<void>;
  answer?: {
    answer: string;
    topicId: string;
  };
  isEditMode?: boolean;
  onCancel?: () => void;
}

const formSchema = z.object({
  type: z.enum([TopicAnswerType.REAL, TopicAnswerType.FAKE]),
  answer: z.string(),
});

export default function TopicMessageAnswer({
  onClickHistory,
  topicId,
  onUpdateAnswer,
  answer,
  isEditMode,
  onCancel,
}: TopicMessageAnswerProps) {
  const { t } = useTranslation();
  const { isLoading, data: answerData, error } = useGetTopicAnswerByTopicId(topicId);
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      answer: answer?.answer || '',
    },
  });

  useEffect(() => {
    if (answerData) {
      form.reset({
        type: answerData.type,
        answer: answerData.answer,
      });
    }
  }, [answerData, form]);

  async function handleSubmit(data: z.infer<typeof formSchema>) {
    if (answer) {
      await onUpdateAnswer(answer.topicId, data.answer);
    }
  }

  if (isLoading) {
    return <LoadingState />;
  }

  if (error) {
    console.log(error);
    return <ErrorState />;
  }

  if (isEditMode) {
    return (
      <Form {...form}>
        <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="answer"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Textarea {...field} />
                </FormControl>
              </FormItem>
            )}
          />
          <div className="flex gap-2">
            <Button type="submit">{t('common.save')}</Button>
            <Button variant="outline" onClick={onCancel}>
              {t('common.cancel')}
            </Button>
          </div>
        </form>
      </Form>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <TYH3>{t('topicMessageAnswer.title')}</TYH3>
        <Button variant="ghost" size="icon" onClick={onClickHistory}>
          <History className="h-4 w-4" />
        </Button>
      </div>
      <div className="rounded-md border p-4">{answer?.answer || t('topicMessageAnswer.noAnswer')}</div>
    </div>
  );
}
