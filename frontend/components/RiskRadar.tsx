import type { ModuleResult } from "@/lib/types";

interface RiskRadarProps {
  modules: string[];
  results: Record<string, ModuleResult>;
}

export function RiskRadar({ modules, results }: RiskRadarProps) {
  const size = 300;
  const center = size / 2;
  const radius = 104;
  const axisRadius = 116;

  const points = modules.map((module, index) => {
    const angle = (Math.PI * 2 * index) / modules.length - Math.PI / 2;
    const score = results[module]?.score ?? 0;
    const scaled = radius * (score / 100);
    return {
      module,
      score,
      x: center + Math.cos(angle) * scaled,
      y: center + Math.sin(angle) * scaled,
      axisX: center + Math.cos(angle) * axisRadius,
      axisY: center + Math.sin(angle) * axisRadius,
      labelX: center + Math.cos(angle) * (axisRadius + 22),
      labelY: center + Math.sin(angle) * (axisRadius + 22)
    };
  });

  const polygon = points.map((point) => `${point.x},${point.y}`).join(" ");

  return (
    <div className="rounded-lg border border-[var(--line)] bg-white p-5 shadow-soft">
      <div className="mb-4 flex items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-semibold">Risk Summary</h2>
          <p className="text-sm text-[var(--muted)]">Module scores on a 0-100 scale</p>
        </div>
      </div>
      <svg viewBox={`0 0 ${size} ${size}`} className="mx-auto h-72 w-full max-w-sm" role="img" aria-label="Module score radar chart">
        {[25, 50, 75, 100].map((level) => {
          const ring = modules
            .map((_, index) => {
              const angle = (Math.PI * 2 * index) / modules.length - Math.PI / 2;
              const scaled = radius * (level / 100);
              return `${center + Math.cos(angle) * scaled},${center + Math.sin(angle) * scaled}`;
            })
            .join(" ");
          return <polygon key={level} points={ring} fill="none" stroke="#dfe4dc" strokeWidth="1" />;
        })}
        {points.map((point) => (
          <line key={point.module} x1={center} y1={center} x2={point.axisX} y2={point.axisY} stroke="#d3d9d0" strokeWidth="1" />
        ))}
        <polygon points={polygon} fill="#10b98133" stroke="#0f9f72" strokeWidth="2" />
        {points.map((point) => (
          <g key={point.module}>
            <circle cx={point.x} cy={point.y} r="4" fill="#0f9f72" />
            <text x={point.labelX} y={point.labelY} textAnchor="middle" dominantBaseline="middle" className="fill-zinc-700 text-[11px] font-semibold">
              {point.module}
            </text>
          </g>
        ))}
      </svg>
    </div>
  );
}
