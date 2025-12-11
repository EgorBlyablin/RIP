import { Alert, Button, Col, Container, Form, Row } from "react-bootstrap";

import { TurbineGenerationRequestCard } from '../components/turbine-card';

import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { deleteGenerationRequest, deleteTurbineFromGenerationRequest, getGenerationRequest, sendGenerationRequest, setGenerationRequestData, updateTurbineInGenerationRequest } from "../slices/generation-requests-draft";
import { useDispatch, useSelector } from '../store';
import { BeatLoader } from "react-spinners";
import { ROUTE_LABELS, ROUTES } from "../routes";
import { Breadcrumbs } from "../components/breadcrumbs";

export const GenerationRequestPage = () => {
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const { generationRequestId } = useParams();

    const { turbines, generationRequestData, loading, isDraft, error } = useSelector(state => state.generationRequestDraft);

    const [periodDaysEntered, setPeriodDaysEntered] = useState<number | undefined>(generationRequestData.period_days)

    useEffect(() => {
        generationRequestId && dispatch(getGenerationRequest(Number(generationRequestId)));
    }, []);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        dispatch(sendGenerationRequest())
            .then((res) => res.meta.requestStatus === "fulfilled" && navigate(ROUTES.TURBINES_LIST))
    };

    const handleDelete = async (e: React.FormEvent) => {
        e.preventDefault();

        dispatch(deleteGenerationRequest())
            .then((res) => res.meta.requestStatus === "fulfilled" && navigate(ROUTES.TURBINES_LIST))
    };

    return (<>
        {loading
            ? <Container style={{ display: "flex", justifyContent: "center", paddingTop: "30%" }}>
                <BeatLoader color="#5BA1D4" />
            </Container>
            : <Container>
                <Breadcrumbs breadcrumbs={[
                    { label: ROUTE_LABELS.HOME, path: ROUTES.HOME },
                    { label: `Заявки`, path: ROUTES.GENERATION_REQUESTS },
                    { label: `№${generationRequestId}`, path: `${ROUTES.GENERATION_REQUESTS}/${generationRequestId}` }
                ]} />
                <h1 style={{ marginBottom: 25 }}>Заявка №{generationRequestId}{generationRequestData.status !== "draft" && <>
                    {" "}({{
                        sent: "В работе",
                        rejected: "Отклонена",
                        completed: "Завершена"
                    }[generationRequestData.status || ""]})
                </>
                }
                </h1>
                {error && <Alert variant="danger">{error}</Alert>}
                <div style={{ fontSize: 18, marginBottom: 30, display: "flex", justifyContent: "space-between" }}>
                    <div>
                        <span style={{ marginRight: 10 }}>Период расчета (дней):</span>
                        {isDraft ?
                            <Form.Control
                                type="number"
                                defaultValue={generationRequestData.period_days}
                                value={periodDaysEntered}
                                onChange={(e) => setPeriodDaysEntered(Number.parseInt(e.currentTarget.value))}
                                style={{ width: 100, display: "inline-block" }}
                                required
                            />
                            : <strong>{generationRequestData.period_days}</strong>
                        }
                        {(periodDaysEntered !== generationRequestData.period_days && !!periodDaysEntered) &&
                            <Button variant="primary" style={{ marginLeft: 10, background: "#5BA1D4", borderColor: "#5BA1D4" }} onClick={() => dispatch(setGenerationRequestData(periodDaysEntered))}>
                                Сохранить
                            </Button>
                        }
                    </div>
                    <div style={{ display: "flex", gap: 10 }}>
                        {(isDraft) && (
                            <>
                                <Button variant="primary" onClick={handleSubmit} style={{ background: "#5BA1D4", borderColor: "#5BA1D4" }}>
                                    Отправить
                                </Button>
                                <Button variant="danger" onClick={handleDelete}>
                                    Очистить
                                </Button>
                            </>
                        )}
                    </div>
                </div>
                <Row style={{
                    fontWeight: 600,
                    padding: "10px 5px",
                    borderBottom: "2px solid #E5E5E5",
                    marginBottom: 10
                }}>
                    <Col md={3} style={{ textAlign: "center" }}>Турбина:</Col>
                    <Col md={1} style={{ textAlign: "center" }}>Мощность:</Col>
                    <Col md={1} style={{ textAlign: "center" }}>Мачта:</Col>
                    <Col md={2} style={{ textAlign: "center" }}>V<sub>ср. ветра</sub>(10 м), м/с:</Col>
                    <Col md={2} style={{ textAlign: "center" }}>α-коэффициент</Col>
                    <Col md={3} style={{ textAlign: "center" }}>Выработка за период</Col>
                </Row>
                {turbines.length
                    ? turbines.map((turbine) => (
                        turbine &&
                        <TurbineGenerationRequestCard
                            turbine={turbine}
                            editable={isDraft}
                            onDelete={() => dispatch(deleteTurbineFromGenerationRequest(turbine.turbine?.id || -1))}
                            onSubmit={(avgWind, alpha) => dispatch(updateTurbineInGenerationRequest({
                                turbineId: turbine.turbine?.id || -1,
                                avgWind: avgWind || undefined,
                                alpha: alpha || undefined
                            }))}
                        />
                    ))
                    : <Row style={{ textAlign: "center", paddingTop: 50 }}>
                        <h3>Заявка пуста</h3>
                    </Row>}
            </Container>
        }
    </>)
};