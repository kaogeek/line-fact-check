interface UserAvatarProps {
  avatarUrl: string;
  name: string;
  className?: string;
}

export default function UserAvatar({ avatarUrl, name, className }: UserAvatarProps) {
  return (
    <figure className={`flex gap-2 items-center ${className ?? ''}`}>
      <img src={avatarUrl} alt={`${name}'s avatar`} className="w-[32px] h-[32px] object-cover rounded-full shadow-md" />
      <figcaption className="text-sm">{name}</figcaption>
    </figure>
  );
}
